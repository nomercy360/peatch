package notification

import (
	"bytes"
	"context"
	"errors"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

type TelegramNotifier struct {
	tg *telegram.Bot
}

func NewTelegramNotifier(bot *telegram.Bot) *TelegramNotifier {
	return &TelegramNotifier{
		tg: bot,
	}
}

type SendNotificationParams struct {
	ChatID    int64
	Message   string
	BotWebApp string
	WebAppURL string
	Image     []byte
	Username  *string
}

func (t *TelegramNotifier) SendNotification(params SendNotificationParams) error {
	log.Printf("Sending notification to chatID: %d", params.ChatID)

	button := models.InlineKeyboardButton{Text: "View"}

	if params.BotWebApp != "" {
		button.URL = params.BotWebApp
	} else if params.WebAppURL != "" {
		button.WebApp = &models.WebAppInfo{URL: params.WebAppURL}
	} else {
		return errors.New("no URL provided")
	}

	buttons := [][]models.InlineKeyboardButton{
		{button},
	}

	if params.Username != nil {
		contact := models.InlineKeyboardButton{Text: "Contact in Telegram", URL: "https://t.me/" + *params.Username}
		buttons = append(buttons, []models.InlineKeyboardButton{contact})
	}

	photoParams := &telegram.SendPhotoParams{
		// ChatID: 927635965,
		ChatID:              params.ChatID,
		Caption:             params.Message,
		ParseMode:           models.ParseModeMarkdown,
		Photo:               &models.InputFileUpload{Filename: "img.jpg", Data: bytes.NewReader(params.Image)},
		DisableNotification: true,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	}

	_, err := t.tg.SendPhoto(context.Background(), photoParams)
	if err != nil {
		log.Printf("Failed to send photo: %s", err)
		return err
	}

	return nil
}
