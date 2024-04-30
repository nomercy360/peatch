package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
	"strconv"
	"time"
)

type UserWithToken struct {
	User      db.User `json:"user"`
	Token     string  `json:"token"`
	Following []int64 `json:"following"`
} // @Name UserWithToken

type TelegramUser struct {
	ID        int64   `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Username  string  `json:"username"`
	Language  *string `json:"language_code"`
	IsPremium bool    `json:"is_premium"`
	AllowsPM  bool    `json:"allows_write_to_pm"`
}

func (tu *TelegramUser) ToDBUser() db.User {
	return db.User{
		FirstName:    tu.FirstName,
		LastName:     tu.LastName,
		Username:     tu.Username,
		LanguageCode: tu.Language,
		ChatID:       tu.ID,
	}
}

func (s *service) TelegramAuth(queryID, userJSON, authDate, hash string) (*UserWithToken, error) {
	if queryID == "" || userJSON == "" || authDate == "" || hash == "" {
		return nil, terrors.BadRequest(errors.New(fmt.Sprintf("missing required query parameters: query_id=%s, user=%s, auth_date=%s, hash=%s", queryID, userJSON, authDate, hash)))
	}

	botToken := s.config.BotToken

	dataCheck := fmt.Sprintf("auth_date=%s\nquery_id=%s\nuser=%s", authDate, queryID, userJSON)

	h1 := hmac.New(sha256.New, []byte("WebAppData"))
	h1.Write([]byte(botToken))
	h2 := hmac.New(sha256.New, h1.Sum(nil))
	h2.Write([]byte(dataCheck))

	computedHash := hex.EncodeToString(h2.Sum(nil))

	if computedHash != hash {
		return nil, errors.New("invalid hash: authentication data may have been tampered with")
	}

	authTimestamp, err := strconv.ParseInt(authDate, 10, 64)
	if err != nil {
		return nil, err
	}
	if time.Since(time.Unix(authTimestamp, 0)) > 24*time.Hour {
		return nil, errors.New("authentication data is outdated")
	}

	var tgUser TelegramUser
	err = json.Unmarshal([]byte(userJSON), &tgUser)
	if err != nil {
		return nil, err
	}

	user, err := s.storage.GetUserByChatID(tgUser.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		user, err = s.storage.CreateUser(tgUser.ToDBUser())
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// fetch following
	following, err := s.storage.GetUserFollowing(user.ID)

	if err != nil {
		return nil, err
	}

	token, err := generateJWT(user.ID, user.ChatID)

	if err != nil {
		return nil, err
	}

	return &UserWithToken{
		User:      *user,
		Token:     token,
		Following: following,
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
