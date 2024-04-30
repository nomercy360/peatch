package notification

import (
	"fmt"
	"log"
	"os"
	"time"
)

type DummyNotifier struct{}

func NewDummyNotifier() *DummyNotifier {
	return &DummyNotifier{}
}

func (t *DummyNotifier) SendNotification(chatID int64, message, imgUrl, link string) error {
	log.Printf("Writing to file %d.json", chatID)

	file, err := os.Create(fmt.Sprintf("%d.json", chatID))
	if err != nil {
		return err
	}

	defer file.Close()

	data := fmt.Sprintf("{\"message\": \"%s\", \"imgUrl\": \"%s\", \"link\": \"%s\"}", message, imgUrl, link)

	if _, err = file.WriteString(data); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	// Simulate some processing time
	time.Sleep(2 * time.Second)

	//log.Printf("Notification successfully sent to chatID: %d with message: '%s', image URL: '%s', and link: '%s'", chatID, message, imgUrl, link)

	return nil
}
