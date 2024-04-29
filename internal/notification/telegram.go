package notification

import "github.com/labstack/gommon/log"

type TelegramNotifier struct {
	BotToken string
}

func (t *TelegramNotifier) SendNotification(chatID int64, message, imgUrl, link string) error {
	log.Infof("Sending notification to chat %d: %s", chatID, message)

	return nil
}
