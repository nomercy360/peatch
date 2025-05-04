package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/nanoid"
)

type LocalizedMessages map[db.LanguageCode]map[string]string

const (
	MsgKeyWelcome        = "welcome"
	MsgKeyOpenWebApp     = "openWebApp"
	MsgKeyLaunch         = "launch"
	MsgKeyOpenWebAppMenu = "openWebAppMenu"
)

var messages = LocalizedMessages{
	db.LanguageEN: {
		MsgKeyWelcome:        "Welcome!\n*Peatch* is a social network made for collaborations. Tap the button to open the web app!",
		MsgKeyOpenWebApp:     "You can open the web app by tapping the button below.",
		MsgKeyLaunch:         "Launch",
		MsgKeyOpenWebAppMenu: "Open Web App",
	},
	db.LanguageRU: {
		MsgKeyWelcome:        "Привет!\n*Peatch* - социальная сеть для совместной работы. Кнопка ниже, откроет веб-приложение!",
		MsgKeyOpenWebApp:     "Вы можете открыть веб-приложение, нажав кнопку ниже.",
		MsgKeyLaunch:         "Запустить",
		MsgKeyOpenWebAppMenu: "Открыть веб-app",
	},
}

func (h *handler) HandleWebhook(c echo.Context) error {
	ctx := c.Request().Context()
	var update models.Update

	if err := json.NewDecoder(c.Request().Body).Decode(&update); err != nil {
		h.logger.Error("failed to decode update", slog.String("error", err.Error()))
		return c.NoContent(http.StatusOK)
	}

	if update.Message == nil {
		return c.NoContent(http.StatusOK)
	}

	if update.Message.Chat.Type != "private" {
		h.logger.Info("ignoring non-private chat", slog.String("chat_type", update.Message.Chat.Type))
		return c.NoContent(http.StatusOK)
	}

	if update.Message.From.IsBot {
		h.logger.Info("ignoring message from bot", slog.String("username", update.Message.From.Username))
		return c.NoContent(http.StatusOK)
	}

	if err := h.handleMessage(ctx, update); err != nil {
		h.logger.Error("handle message failed", slog.String("error", err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) handleMessage(ctx context.Context, update models.Update) error {
	chatID := update.Message.Chat.ID

	lang := determineLanguage(update.Message.From.LanguageCode)
	msgs := messages[lang]

	webAppKeyboard := createWebAppKeyboard(msgs[MsgKeyLaunch], h.config.WebAppURL)

	user, err := h.storage.GetUserByChatID(ctx, chatID)

	if errors.Is(err, db.ErrNotFound) {

		h.logger.Info("creating new user",
			slog.String("chat_id", fmt.Sprintf("%d", chatID)),
			slog.String("username", update.Message.From.Username))

		user = h.createUser(ctx, update, lang)
		if user.ID == "" {
			h.logger.Error("failed to create user", slog.Int64("chat_id", chatID))
			return errors.New("failed to create user")
		}

		if err := h.sendWelcomePhoto(ctx, chatID, msgs[MsgKeyWelcome], webAppKeyboard); err != nil {
			h.logger.Error("failed to send welcome photo",
				slog.Int64("chat_id", chatID),
				slog.String("error", err.Error()))
			return err
		}

		go h.setMenuButton(chatID, lang)

	} else if err != nil {
		h.logger.Error("failed to query user",
			slog.Int64("chat_id", chatID),
			slog.String("error", err.Error()))
		return err
	} else {

		h.logger.Info("existing user interaction",
			slog.Int64("chat_id", chatID),
			slog.String("user_id", user.ID))

		params := &telegram.SendMessageParams{
			ChatID:      chatID,
			Text:        msgs[MsgKeyOpenWebApp],
			ReplyMarkup: webAppKeyboard,
			ParseMode:   "Markdown",
		}

		if _, err := h.bot.SendMessage(ctx, params); err != nil {
			h.logger.Error("failed to send message",
				slog.Int64("chat_id", chatID),
				slog.String("error", err.Error()))
			return err
		}
	}

	return nil
}

func determineLanguage(langCode string) db.LanguageCode {
	if langCode == "ru" {
		return db.LanguageRU
	}
	return db.LanguageEN
}

func createWebAppKeyboard(buttonText, webAppURL string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:   buttonText,
					WebApp: &models.WebAppInfo{URL: webAppURL},
				},
			},
		},
	}
}

func (h *handler) sendWelcomePhoto(ctx context.Context, chatID int64, message string, keyboard *models.InlineKeyboardMarkup) error {
	photo := &models.InputFileString{Data: "https://assets.peatch.io/peatch-preview.png"}

	params := &telegram.SendPhotoParams{
		ChatID:      chatID,
		Caption:     message,
		ReplyMarkup: keyboard,
		Photo:       photo,
		ParseMode:   "Markdown",
	}

	_, err := h.bot.SendPhoto(ctx, params)
	return err
}

func (h *handler) handleBotAvatar(ctx context.Context, userID string, chatID int64) {
	photos, err := h.bot.GetUserProfilePhotos(ctx, &telegram.GetUserProfilePhotosParams{
		UserID: chatID,
		Offset: 0,
		Limit:  1,
	})
	if err != nil {
		h.logger.Error("failed to get profile photos",
			slog.Int64("chat_id", chatID),
			slog.String("error", err.Error()))
		return
	}

	if photos.TotalCount == 0 {
		h.logger.Info("user has no profile photos", slog.Int64("chat_id", chatID))
		return
	}

	bestPhoto := findBestQualityPhoto(photos.Photos)
	if bestPhoto == nil {
		h.logger.Error("could not find suitable photo", slog.Int64("chat_id", chatID))
		return
	}

	file, err := h.bot.GetFile(ctx, &telegram.GetFileParams{FileID: bestPhoto.FileID})
	if err != nil {
		h.logger.Error("failed to get file",
			slog.String("file_id", bestPhoto.FileID),
			slog.String("error", err.Error()))
		return
	}

	if err := h.downloadAndStoreAvatar(ctx, userID, chatID, h.bot.FileDownloadLink(file)); err != nil {
		h.logger.Error("avatar processing failed",
			slog.Int64("chat_id", chatID),
			slog.String("error", err.Error()))
		return
	}

	h.logger.Info("avatar processed successfully", slog.Int64("chat_id", chatID))
}

func findBestQualityPhoto(photos [][]models.PhotoSize) *models.PhotoSize {
	var bestPhoto *models.PhotoSize

	for _, album := range photos {
		for i := range album {
			photo := &album[i]
			if bestPhoto == nil ||
				photo.FileSize > bestPhoto.FileSize ||
				(photo.FileSize == bestPhoto.FileSize && photo.Width > bestPhoto.Width) {
				bestPhoto = photo
			}
		}
	}

	return bestPhoto
}

func (h *handler) downloadAndStoreAvatar(ctx context.Context, userID string, chatID int64, fileURL string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	fileName := fmt.Sprintf("%s/%s.jpg", userID, nanoid.Must())

	if err = h.s3Client.UploadFile(ctx, fileName, resp.Body, "image/jpeg"); err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	if err := h.storage.UpdateUserAvatarURL(ctx, userID, fileName); err != nil {
		return fmt.Errorf("failed to update avatar URL: %w", err)
	}

	return nil
}

func (h *handler) setMenuButton(chatID int64, lang db.LanguageCode) {
	ctx := context.Background()
	msg := messages[lang]

	menu := telegram.SetChatMenuButtonParams{
		ChatID: chatID,
		MenuButton: models.MenuButtonWebApp{
			Type:   "web_app",
			Text:   msg[MsgKeyOpenWebAppMenu],
			WebApp: models.WebAppInfo{URL: h.config.WebAppURL},
		},
	}

	if _, err := h.bot.SetChatMenuButton(ctx, &menu); err != nil {
		h.logger.Error("failed to set menu button",
			slog.Int64("chat_id", chatID),
			slog.String("error", err.Error()))
		return
	}

	h.logger.Info("menu button set successfully", slog.Int64("chat_id", chatID))
}

func (h *handler) createUser(ctx context.Context, update models.Update, lang db.LanguageCode) db.User {
	chatID := update.Message.Chat.ID
	var firstName, lastName *string

	if update.Message.Chat.FirstName != "" {
		firstName = &update.Message.Chat.FirstName
	}

	if update.Message.Chat.LastName != "" {
		lastName = &update.Message.Chat.LastName
	}

	username := update.Message.Chat.Username
	if username == "" {
		username = fmt.Sprintf("user_%d", chatID)
	}

	user := db.User{
		ID:           nanoid.Must(),
		ChatID:       chatID,
		FirstName:    firstName,
		LastName:     lastName,
		Username:     username,
		LanguageCode: lang,
	}

	if err := h.storage.CreateUser(ctx, user); err != nil {
		h.logger.Error("failed to create user",
			slog.Int64("chat_id", chatID),
			slog.String("error", err.Error()))
		return db.User{}
	}

	newUser, err := h.storage.GetUserByChatID(ctx, chatID)
	if err != nil {
		h.logger.Error("failed to retrieve new user",
			slog.Int64("chat_id", chatID),
			slog.String("error", err.Error()))
		return user
	}

	go h.handleBotAvatar(context.Background(), newUser.ID, chatID)

	return newUser
}
