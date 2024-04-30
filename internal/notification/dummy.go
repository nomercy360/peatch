package notification

import (
	"log"
	"time"
)

type DummyNotifier struct{}

func NewDummyNotifier() *DummyNotifier {
	return &DummyNotifier{}
}

func (t *DummyNotifier) SendNotification(chatID int64, message, imgUrl, link string) error {
	log.Printf("Mock Sending notification to chatID: %d", chatID)
	log.Printf("Message: %s", message)
	log.Printf("Image URL: %s", imgUrl)
	log.Printf("Link: %s", link)

	// Simulate some processing time
	time.Sleep(2 * time.Second)

	log.Printf("Notification successfully sent to chatID: %d with message: '%s', image URL: '%s', and link: '%s'", chatID, message, imgUrl, link)

	return nil
}
