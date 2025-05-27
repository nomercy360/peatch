package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"net/http"
	"strings"
)

// AdminAuthGetter is a function type for getting admin by API token
type AdminAuthGetter func(ctx context.Context, apiToken string) (adminID string, err error)

// AdminAuth creates a middleware that supports both JWT and API token authentication
func AdminAuth(jwtSecret string, getAdminByToken AdminAuthGetter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// First, try API token authentication
			apiToken := c.Request().Header.Get("X-API-Token")
			if apiToken != "" {
				adminID, err := getAdminByToken(c.Request().Context(), apiToken)
				if err == nil && adminID != "" {
					// Create a token with claims for consistency
					claims := &contract.AdminJWTClaims{
						AdminID: adminID,
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					c.Set("user", token)
					return next(c)
				}
			}

			// If API token auth failed or not provided, try JWT auth
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")

				token, err := jwt.ParseWithClaims(tokenString, &contract.AdminJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("invalid signing method")
					}
					return []byte(jwtSecret), nil
				})

				if err == nil && token.Valid {
					if claims, ok := token.Claims.(*contract.AdminJWTClaims); ok && claims.AdminID != "" {
						c.Set("user", token)
						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusUnauthorized, ErrAuthInvalid)
		}
	}
}
