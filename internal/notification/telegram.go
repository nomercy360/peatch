package notification

import "github.com/labstack/gommon/log"

type TelegramNotifier struct {
	BotToken string
}

func (t *TelegramNotifier) SendNotification(userID int64, message string) error {
	log.Infof("Sending notification to user %d: %s", userID, message)
	return nil
}
