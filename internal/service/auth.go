package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
	"github.com/telegram-mini-apps/init-data-golang"
	"regexp"
	"strings"
	"time"
)

type UserWithToken struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
} // @Name UserWithToken

func (s *service) TelegramAuth(query string) (*UserWithToken, error) {
	expIn := 24 * time.Hour
	botToken := s.config.BotToken

	if err := initdata.Validate(query, botToken, expIn); err != nil {
		return nil, terrors.Unauthorized(err, "invalid data")
	}

	data, err := initdata.Parse(query)

	if err != nil {
		return nil, terrors.Unauthorized(err, "invalid data")
	}

	user, err := s.storage.GetUserByChatID(data.User.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		var firstName, lastName, langCode *string

		if data.User.FirstName != "" {
			firstName = &data.User.FirstName
		}

		if data.User.LastName != "" {
			lastName = &data.User.LastName
		}

		if data.User.LanguageCode != "" {
			langCode = &data.User.LanguageCode
		}

		username := data.User.Username
		if username == "" && firstName != nil {
			username = urlify(*firstName)
		} else if username == "" {
			username = "user" + fmt.Sprintf("%d", data.User.ID)
		}

		create := db.User{
			FirstName:    firstName,
			LastName:     lastName,
			Username:     data.User.Username,
			LanguageCode: langCode,
			ChatID:       data.User.ID,
		}

		user, err = s.storage.CreateUser(create)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	token, err := generateJWT(user.ID, user.ChatID)

	if err != nil {
		return nil, err
	}

	return &UserWithToken{
		User:  *user,
		Token: token,
	}, nil
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UID    int64 `json:"uid"`
	ChatID int64 `json:"chat_id"`
}

func generateJWT(id, chatID int64) (string, error) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID:    id,
		ChatID: chatID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}

func urlify(s string) string {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, " ", "_")

	reg := regexp.MustCompile(`[^a-z0-9_]+`)
	s = reg.ReplaceAllString(s, "_")

	s = strings.Trim(s, "_")

	return s
}
