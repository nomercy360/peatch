package notification

import (
	"bytes"
	"context"
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

func (t *TelegramNotifier) SendNotification(chatID int64, message, link string, img []byte) error {
	log.Printf("Sending notification to chatID: %d", chatID)

	photoParams := &telegram.SendPhotoParams{
		//ChatID:              927635965,
		ChatID:              chatID,
		Caption:             message,
		ParseMode:           models.ParseModeMarkdown,
		Photo:               &models.InputFileUpload{Filename: "img.jpg", Data: bytes.NewReader(img)},
		DisableNotification: true,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "View", URL: "t.me/peatcher_testing_bot/peatch"},
				},
			},
		},
	}

	_, err := t.tg.SendPhoto(context.Background(), photoParams)
	if err != nil {
		log.Printf("Failed to send photo: %s", err)
		return err
	}

	return nil
}
