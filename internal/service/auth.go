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
	botID := s.config.BotID // bot id is a space for future use

	if err := initdata.ValidateThirdParty(query, botID, expIn); err != nil {
		return nil, terrors.Unauthorized(err, "invalid data")
	}

	data, err := initdata.Parse(query)

	if err != nil {
		return nil, terrors.Unauthorized(err, "invalid data")
	}

	user, err := s.storage.GetUserByChatID(data.User.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		var firstName, lastName *string

		if data.User.FirstName != "" {
			firstName = &data.User.FirstName
		}

		if data.User.LastName != "" {
			lastName = &data.User.LastName
		}

		username := data.User.Username
		if username == "" {
			username = "user_" + fmt.Sprintf("%d", data.User.ID)
		}

		create := db.User{
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			ChatID:    data.User.ID,
		}

		lang := "ru"

		if data.User.LanguageCode != "ru" {
			lang = "en"
		}

		create.LanguageCode = lang

		user, err = s.storage.CreateUser(create)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	token, err := generateJWT(user.ID, user.ChatID, user.LanguageCode)

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
	UID    int64  `json:"uid"`
	ChatID int64  `json:"chat_id"`
	Lang   string `json:"lang"`
}

func generateJWT(id, chatID int64, lang string) (string, error) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID:    id,
		ChatID: chatID,
		Lang:   lang,
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
