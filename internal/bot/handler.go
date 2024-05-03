package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	telegram "github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
	"github.com/peatch-io/peatch/internal/db"
	"io"
	"log"
	"net/http"
)

type bot struct {
	storage  storage
	config   Config
	s3Client s3Client
	tg       *telegram.Bot
	messages map[string]map[string]string
}

type storage interface {
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) (*db.User, error)
	UpdateUserAvatarURL(chatID int64, avatarURL string) error
	DeleteUserByID(id int64) error
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
		messages: loadTranslations(),
	}
	b.initTelegram()
	return b
}

func loadTranslations() map[string]map[string]string {
	return map[string]map[string]string{
		"en": {
			"welcome":        "Welcome!\n*Peatch* is a social network made for collaborations. Tap the button to open the web app!",
			"openWebApp":     "You can open the web app by tapping the button below.",
			"launch":         "Launch",
			"openWebAppMenu": "Open Web App",
		},
		"ru": {
			"welcome":        "Привет!\n*Peatch* - социальная сеть для совместной работы. Кнопка ниже, откроет веб-приложение!",
			"openWebApp":     "Вы можете открыть веб-приложение, нажав кнопку ниже.",
			"launch":         "Запустить",
			"openWebAppMenu": "Открыть веб-app",
		},
	}
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
	// only respond to messages from users. Ignore messages from bots and groups. ALso ignore messages from channels
	if update.Message.Chat.Type != "private" {
		w.WriteHeader(http.StatusOK)
		return
	} else if update.Message.From.IsBot {
		w.WriteHeader(http.StatusOK)
		return
	}

	user, err := b.storage.GetUserByChatID(update.Message.Chat.ID)

	// if its /reset command
	if update.Message.Text == "/reset" && user != nil {
		if err := b.storage.DeleteUserByID(user.ID); err != nil {
			log.Printf("Failed to delete user: %v", err)
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		msg := telegram.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "User deleted",
		}

		if _, err := b.tg.SendMessage(context.Background(), &msg); err != nil {
			log.Printf("Failed to send message: %v", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	lang := "en"
	if update.Message.From.LanguageCode != "" {
		lang = update.Message.From.LanguageCode
	}

	msgs := b.messages[lang]

	webApp := tgModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgModels.InlineKeyboardButton{
			{
				{Text: msgs["launch"], WebApp: &tgModels.WebAppInfo{URL: b.config.WebAppURL}},
			},
		},
	}

	if err != nil && errors.Is(err, db.ErrNotFound) {
		log.Printf("User %d not found, creating new user", update.Message.Chat.ID)

		user = b.createUser(update)
		if user == nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		message := msgs["welcome"]

		photo := &tgModels.InputFileString{Data: "https://assets.peatch.io/peatch-preview.png"}

		params := &telegram.SendPhotoParams{ChatID: update.Message.Chat.ID, Caption: message, ReplyMarkup: &webApp, Photo: photo, ParseMode: "Markdown"}

		if _, err := b.tg.SendPhoto(context.Background(), params); err != nil {
			log.Printf("Failed to send message: %v", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		go b.setMenuButton(update.Message.Chat.ID, lang)

	} else if err != nil {
		log.Printf("Failed to get user: %v", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	} else {
		log.Printf("User %d already exists, sending message", user.ChatID)

		message := msgs["openWebApp"]

		params := &telegram.SendMessageParams{ChatID: update.Message.Chat.ID, Text: message, ReplyMarkup: &webApp, ParseMode: "Markdown"}

		if _, err := b.tg.SendMessage(context.Background(), params); err != nil {
			log.Printf("Failed to send message: %v", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (b *bot) createUser(update tgModels.Update) *db.User {
	// Extract user details from update
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

	go b.handleUserAvatar(newUser.ID, update.Message.From.ID, newUser.ChatID)

	return newUser
}

func (b *bot) handleUserAvatar(userID, tgUserID, chatID int64) {
	photos, err := b.tg.GetUserProfilePhotos(context.Background(), &telegram.GetUserProfilePhotosParams{UserID: tgUserID, Offset: 0, Limit: 1})
	if err != nil {
		log.Printf("Failed to get user profile photos: %v", err)
		return
	}

	if photos.TotalCount > 0 {
		photo := photos.Photos[0][0]

		file, err := b.tg.GetFile(context.Background(), &telegram.GetFileParams{FileID: photo.FileID})
		if err != nil {
			log.Printf("Failed to get file: %v", err)
			return
		}

		fileURL := b.tg.FileDownloadLink(file)

		resp, err := http.Get(fileURL)

		if err != nil {
			log.Printf("Failed to download file: %v", err)
			return
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Failed to read file: %v", err)
			return
		}

		fileName := fmt.Sprintf("%d/%d.jpg", userID, chatID)

		if _, err = b.s3Client.UploadFile(data, fileName); err != nil {
			log.Printf("Failed to upload user avatar to S3: %v", err)
			return
		}

		log.Printf("Avatar uploaded successfully: %s", fileName)

		if err := b.storage.UpdateUserAvatarURL(userID, fileName); err != nil {
			log.Printf("Failed to update user avatar URL: %v", err)
		}

		log.Printf("Profile photo updated for user %d", chatID)
	}
}

func (b *bot) SetWebhook() error {
	webhook := &telegram.SetWebhookParams{URL: b.config.ExternalURL + "/webhook"}

	if _, err := b.tg.SetWebhook(context.Background(), webhook); err != nil {
		return err
	}

	return nil
}

func (b *bot) setMenuButton(chatID int64, lang string) {
	msg := b.messages[lang]

	menu := telegram.SetChatMenuButtonParams{
		ChatID: chatID,
		MenuButton: tgModels.MenuButtonWebApp{
			Type:   "web_app",
			Text:   msg["openWebAppMenu"],
			WebApp: tgModels.WebAppInfo{URL: b.config.WebAppURL},
		},
	}

	if _, err := b.tg.SetChatMenuButton(context.Background(), &menu); err != nil {
		log.Printf("Failed to set chat menu button: %v", err)
		return
	}

	log.Printf("User %d menu button set", chatID)
}
