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
	"regexp"
	"strconv"
	"strings"
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
	GetUserProfile(params db.GetUsersParams) (*db.User, error)
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
		log.Printf("Failed to decode update: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	if update.Message == nil {
		w.WriteHeader(http.StatusOK)
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
			w.WriteHeader(http.StatusOK)
			return
		}

		msg := telegram.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "User deleted",
		}

		if _, err := b.tg.SendMessage(context.Background(), &msg); err != nil {
			log.Printf("Failed to send message: %v", err)
			w.WriteHeader(http.StatusOK)
			return
		}

		return
	}

	lang := "ru"
	if update.Message.From.LanguageCode != "ru" {
		lang = "en"
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
			w.WriteHeader(http.StatusOK)
			return
		}

		message := msgs["welcome"]

		photo := &tgModels.InputFileString{Data: "https://assets.peatch.io/peatch-preview.png"}

		params := &telegram.SendPhotoParams{ChatID: update.Message.Chat.ID, Caption: message, ReplyMarkup: &webApp, Photo: photo, ParseMode: "Markdown"}

		if _, err := b.tg.SendPhoto(context.Background(), params); err != nil {
			log.Printf("Failed to send message: %v", err)
			w.WriteHeader(http.StatusOK)
			return
		}

		go b.setMenuButton(update.Message.Chat.ID, lang)

	} else if err != nil {
		log.Printf("Failed to get user: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		log.Printf("User %d already exists, sending message", user.ChatID)

		message := msgs["openWebApp"]

		params := &telegram.SendMessageParams{ChatID: update.Message.Chat.ID, Text: message, ReplyMarkup: &webApp, ParseMode: "Markdown"}

		if _, err := b.tg.SendMessage(context.Background(), params); err != nil {
			w.WriteHeader(http.StatusOK)
			log.Printf("Failed to send message: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func extractReferrerID(arg string) int64 {
	idStr := strings.TrimPrefix(arg, "friend")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Failed to parse referrer ID: %v", err)
		return 0
	}
	return id
}

func (b *bot) createUser(update tgModels.Update) *db.User {
	// Extract user details from update
	var firstName, lastName *string
	if update.Message.Chat.FirstName != "" {
		firstName = &update.Message.Chat.FirstName
	}

	if update.Message.Chat.LastName != "" {
		lastName = &update.Message.Chat.LastName
	}

	// if username is empty, use first name
	username := update.Message.Chat.Username

	if username == "" {
		username = "user_" + fmt.Sprintf("%d", update.Message.Chat.ID)
	}

	user := db.User{
		ChatID:    update.Message.Chat.ID,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	}

	lang := "ru"

	if update.Message.From.LanguageCode != "ru" {
		lang = "en"
	}

	user.LanguageCode = &lang

	if strings.HasPrefix(update.Message.Text, "/start") {
		// /start friend123 or can be just /start
		args := strings.Fields(update.Message.Text)
		if len(args) > 1 && strings.HasPrefix(args[1], "friend") {
			id := extractReferrerID(args[1])
			if id > 0 {
				if _, err := b.storage.GetUserProfile(db.GetUsersParams{UserID: id, ViewerID: id}); err != nil {
					log.Printf("Failed to get referrer: %v", err)
				} else {
					user.ReferrerID = &id
				}
			}
		}
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
		bestPhoto := new(tgModels.PhotoSize)

		for _, album := range photos.Photos {
			for _, pic := range album {
				if pic.FileSize > bestPhoto.FileSize || (pic.FileSize == bestPhoto.FileSize && pic.Width > bestPhoto.Width) {
					bestPhoto = &pic
				}
			}
		}

		file, err := b.tg.GetFile(context.Background(), &telegram.GetFileParams{FileID: bestPhoto.FileID})
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
	webhook := &telegram.SetWebhookParams{URL: b.config.ExternalURL + "/webhook", MaxConnections: 10}

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

func urlify(s string) string {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, " ", "_")

	reg := regexp.MustCompile(`[^a-z0-9_]+`)
	s = reg.ReplaceAllString(s, "_")

	s = strings.Trim(s, "_")

	return s
}
