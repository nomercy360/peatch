package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
	"strconv"
	"time"
)

type UserWithToken struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
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

	var userInfo db.User
	err = json.Unmarshal([]byte(userJSON), &userInfo)
	if err != nil {
		return nil, err
	}

	user, err := s.storage.GetUserByChatID(userInfo.ChatID)
	if err != nil {
		if db.IsNoRowsError(err) {
			user, err = s.storage.CreateUser(userInfo)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
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

func generateJWT(id, chatID int64) (string, error) {
	mySigningKey := []byte("secret")

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Issuer:    strconv.FormatInt(chatID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return ss, nil
}
