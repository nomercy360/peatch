package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	initdata "github.com/telegram-mini-apps/init-data-golang"

	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/nanoid"
)

const (
	ErrInvalidInitData = "invalid init data from telegram"
	ErrInvalidRequest  = "failed to validate request"
	ErrAuthInvalid     = "auth is invalid"
)

// TelegramAuth godoc
// @Summary Telegram Auth
// @Description Authenticate user via Telegram using init data
// @Tags auth
// @Param request body contract.AuthTelegramRequest true "Telegram Auth Request"
// @UserStatus 200 {object} contract.AuthResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 500 {object} contract.ErrorResponse
// @Router /auth-telegram [post]
func (h *handler) TelegramAuth(c echo.Context) error {
	var req contract.AuthTelegramRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	h.logger.Info("telegram auth request", slog.String("request", fmt.Sprintf("%+v", req.Query)))

	expIn := 24 * time.Hour
	botToken := h.config.TelegramBotToken

	if err := initdata.Validate(req.Query, botToken, expIn); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidInitData).WithInternal(err)
	}

	data, err := initdata.Parse(req.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidInitData).WithInternal(err)
	}

	user, err := h.storage.GetUserByChatID(c.Request().Context(), data.User.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		username := data.User.Username
		if username == "" {
			username = "u" + fmt.Sprintf("%d", data.User.ID)
		}

		var first, last *string
		if data.User.FirstName != "" {
			first = &data.User.FirstName
		}
		if data.User.LastName != "" {
			last = &data.User.LastName
		}

		lang := db.LanguageRU
		if data.User.LanguageCode != string(db.LanguageRU) {
			lang = db.LanguageEN
		}

		create := db.User{
			ID:                 nanoid.Must(),
			Username:           username,
			ChatID:             data.User.ID,
			FirstName:          first,
			LastName:           last,
			LanguageCode:       lang,
			VerificationStatus: db.VerificationStatusPending,
		}

		if err = h.storage.CreateUser(c.Request().Context(), create); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user").WithInternal(err)
		}

		user, err = h.storage.GetUserByChatID(c.Request().Context(), data.User.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
		}
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	cfCountry := c.Request().Header.Get("CF-IPCountry")
	cfCity := c.Request().Header.Get("CF-IPCity")
	ua := c.Request().UserAgent()

	if err := h.storage.UpdateUserLoginMetadata(c.Request().Context(), user.ID, db.LoginMeta{
		Country:   cfCountry,
		City:      cfCity,
		UserAgent: ua,
	}); err != nil {
		h.logger.Warn("failed to update login metadata", slog.String("user_id", user.ID), slog.Any("error", err))
	}

	token, err := generateJWT(user.ID, user.ChatID, string(user.LanguageCode), h.config.JWTSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "jwt library error").WithInternal(err)
	}

	resp := contract.AuthResponse{
		Token: token,
		User:  user,
	}

	return c.JSON(http.StatusOK, resp)
}

func generateJWT(userID string, telegramID int64, lang string, secretKey string) (string, error) {
	claims := &contract.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID:    userID,
		ChatID: telegramID,
		Lang:   lang,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}
