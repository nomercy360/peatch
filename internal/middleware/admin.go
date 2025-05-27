package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"net/http"
	"strings"
)

func GetAdminAuthConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return new(contract.AdminJWTClaims)
		},
		SigningKey:             []byte(secret),
		ContinueOnIgnoredError: true,
		ErrorHandler: func(c echo.Context, err error) error {

			var extErr *echojwt.TokenExtractionError
			if !errors.As(err, &extErr) {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrAuthInvalid)
			}

			claims := &contract.AdminJWTClaims{}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			c.Set("user", token)

			if claims.AdminID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrAuthInvalid)
			}

			return nil
		},
	}
}

type AdminAPITokenAuth struct {
	storage adminStorage
}

type adminStorage interface {
	GetAdminByAPIToken(ctx context.Context, apiToken string) (db.Admin, error)
}

func NewAdminAPITokenAuth(storage adminStorage) *AdminAPITokenAuth {
	return &AdminAPITokenAuth{storage: storage}
}

// Middleware validates API token from Authorization header
func (a *AdminAPITokenAuth) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get authorization header
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization token")
			}

			// Check for Bearer token format
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization format")
			}

			token := parts[1]

			// Validate API token
			admin, err := a.storage.GetAdminByAPIToken(context.Background(), token)
			if err != nil {
				if errors.Is(err, db.ErrNotFound) {
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization token")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to validate token")
			}

			// Create JWT claims for compatibility with existing code
			claims := &contract.AdminJWTClaims{
				AdminID: admin.ID,
			}
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Store in context
			c.Set("user", jwtToken)

			return next(c)
		}
	}
}
