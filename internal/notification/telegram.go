package notification

import (
	"bytes"
	"context"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io"
	"log"
	"net/http"
	"time"
)

type TelegramNotifier struct {
	tg *telegram.Bot
}

func NewTelegramNotifier(bot *telegram.Bot) *TelegramNotifier {
	return &TelegramNotifier{
		tg: bot,
	}
}

func (t *TelegramNotifier) SendNotification(chatID int64, message, imgUrl, link string) error {
	log.Printf("Sending notification to chatID: %d", chatID)

	httpClient := http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := httpClient.Get(imgUrl)
	if err != nil {
		log.Printf("Failed to download image: %s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download image, got status code: %d", resp.StatusCode)
		return err
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read image data: %s", err)
		return err
	}

	photoParams := &telegram.SendPhotoParams{
		//ChatID:              chatID,
		ChatID:              927635965,
		Caption:             message,
		Photo:               &models.InputFileUpload{Filename: "img.jpg", Data: bytes.NewReader(imgData)},
		DisableNotification: true,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "View", WebApp: &models.WebAppInfo{URL: link}},
				},
			},
		},
	}

	_, err = t.tg.SendPhoto(context.Background(), photoParams)
	if err != nil {
		log.Printf("Failed to send photo: %s", err)
		return err
	}

	return nil
}
