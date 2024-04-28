package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	telegram "github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
	"github.com/peatch-io/peatch/internal/db"
	"log"
	"net/http"
)

type bot struct {
	storage  storage
	config   Config
	s3Client s3Client
	tg       *telegram.Bot
}

type storage interface {
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) (*db.User, error)
}

type s3Client interface {
	UploadFile(file []byte, fileName string) (string, error)
}

type Config struct {
	BotToken    string
	WebAppURL   string
	ExternalURL string
}

func New(s storage, s3 s3Client, config Config) *bot {
	b := &bot{
		storage:  s,
		s3Client: s3,
		config:   config,
	}
	b.initTelegram()
	return b
}

func (b *bot) initTelegram() {
	var err error
	b.tg, err = telegram.New(b.config.BotToken)
	if err != nil {
		log.Fatalf("Failed to create telegram bot: %v", err)
	}
}

func (b *bot) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var update tgModels.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid update", http.StatusBadRequest)
		return
	}

	if update.Message == nil {
		http.Error(w, "No message", http.StatusBadRequest)
		return
	}

	b.handleMessage(update, w)
}

func (b *bot) handleMessage(update tgModels.Update, w http.ResponseWriter) {
	user, err := b.storage.GetUserByChatID(update.Message.Chat.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		log.Printf("User not found, creating new user")

		user = b.createUser(update)
		if user == nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		log.Printf("Failed to get user: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Hello, %s!", *user.FirstName)

	webApp := tgModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgModels.InlineKeyboardButton{
			{
				{Text: "Open Me", WebApp: &tgModels.WebAppInfo{URL: b.config.WebAppURL}},
			},
		},
	}

	params := &telegram.SendMessageParams{ChatID: update.Message.Chat.ID, Text: message, ReplyMarkup: &webApp}

	if _, err := b.tg.SendMessage(context.Background(), params); err != nil {
		log.Printf("Failed to send message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (b *bot) createUser(update tgModels.Update) *db.User {
	var firstName, lastName, langCode *string
	if update.Message.Chat.FirstName != "" {
		firstName = &update.Message.Chat.FirstName
	}

	if update.Message.Chat.LastName != "" {
		lastName = &update.Message.Chat.LastName
	}

	if update.Message.From.LanguageCode != "" {
		langCode = &update.Message.From.LanguageCode
	}

	user := db.User{
		ChatID:       update.Message.Chat.ID,
		Username:     update.Message.Chat.Username,
		FirstName:    firstName,
		LastName:     lastName,
		LanguageCode: langCode,
	}

	newUser, err := b.storage.CreateUser(user)

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil
	}

	return newUser
}

func (b *bot) SetWebhook() error {
	webhook := &telegram.SetWebhookParams{URL: b.config.ExternalURL + "/webhook"}

	if _, err := b.tg.SetWebhook(context.Background(), webhook); err != nil {
		return err
	}

	return nil
}
