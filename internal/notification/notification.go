package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/peatch-io/peatch/internal/db"
)

type NotifierConfig struct {
	BotToken        string
	AdminChatID     int64
	CommunityChatID int64
	BotWebApp       string
	WebAppURL       string
	AdminWebApp     string
	ImageServiceURL string
}

type Notifier struct {
	bot             *telegram.Bot
	adminChatID     int64
	communityChatID int64
	botWebApp       string
	webappURL       string
	adminWebApp     string
	imageServiceURL string
}

func NewNotifier(config NotifierConfig) (*Notifier, error) {
	bot, err := telegram.New(config.BotToken)
	if err != nil {
		return nil, err
	}

	return &Notifier{
		bot:             bot,
		adminChatID:     config.AdminChatID,
		communityChatID: config.CommunityChatID,
		botWebApp:       config.BotWebApp,
		webappURL:       config.WebAppURL,
		adminWebApp:     config.AdminWebApp,
		imageServiceURL: config.ImageServiceURL,
	}, nil
}

type ImageRequest struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Avatar   string `json:"avatar"`
	Tags     []Tag  `json:"tags"`
}

type Tag struct {
	Text  string `json:"text"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func (n *Notifier) NotifyUserVerified(user db.User) error {
	var msgText string
	if user.LanguageCode == db.LanguageRU {
		msgText = fmt.Sprintf("🎉 Поздравляем! Ваш профиль был подтверждён.\nТеперь вы можете пользоваться всеми функциями и быть видимы другим пользователям.\n\nПроверьте свой профиль")
	} else {
		msgText = fmt.Sprintf("🎉 Congratulations! Your profile has been verified.\nYou can now access all features and be visible to other users.\n\nCheck your profile")
	}

	button := models.InlineKeyboardButton{
		Text: "View Profile",
		URL:  fmt.Sprintf("%s?startapp=u_%s", n.botWebApp, user.ID),
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		//ChatID: fmt.Sprintf("%d", user.ChatID),
		ChatID:      n.adminChatID,
		Text:        msgText,
		ReplyMarkup: &keyboard,
	})

	if err != nil {
		return err
	}

	if !user.IsProfileComplete() {
		log.Printf("User %s has an incomplete profile. Not sending welcome message.", user.Username)
		return nil
	}

	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}

	fullName := fmt.Sprintf("%s %s", firstName, lastName)
	if firstName == "" && lastName == "" {
		fullName = user.Username
	}

	tags := make([]Tag, 0, 5)

	if user.Badges != nil && len(user.Badges) > 0 {
		for i, badge := range user.Badges {
			if i >= 5 {
				break
			}

			tag := Tag{
				Text:  badge.Text,
				Color: badge.Color,
				Icon:  badge.Icon,
			}

			tags = append(tags, tag)
		}
	}

	imageReq := ImageRequest{
		Title: fullName,
		Tags:  tags,
	}

	if user.Title != nil {
		imageReq.Subtitle = *user.Title
	}

	if user.AvatarURL != nil {
		imageReq.Avatar = fmt.Sprintf("https://assets.peatch.io/%s", *user.AvatarURL)
	}

	var imageBytes []byte
	if n.imageServiceURL != "" {
		imageBytes, err = n.generateProfileImage(imageReq)
		if err != nil {
			fmt.Printf("Error generating profile image: %v\n", err)
		}
	}

	communityMsg := fmt.Sprintf("🌟 Welcome new member!\nMeet %s\n\nCheck their profile",
		fullName)

	if imageBytes != nil && len(imageBytes) > 0 {
		photoData := &models.InputFileUpload{
			Filename: fmt.Sprintf("welcome_%s.png", user.ID),
			Data:     bytes.NewReader(imageBytes),
		}

		if _, err := n.bot.SendPhoto(context.Background(), &telegram.SendPhotoParams{
			ChatID:      fmt.Sprintf("%d", n.communityChatID),
			Caption:     communityMsg,
			Photo:       photoData,
			ReplyMarkup: &keyboard,
		}); err != nil {
			fmt.Printf("Error sending photo: %v\n", err)
		}
	} else {
		params := &telegram.SendMessageParams{
			ChatID:      fmt.Sprintf("%d", n.communityChatID),
			Text:        communityMsg,
			ReplyMarkup: &keyboard,
		}

		if _, err := n.bot.SendMessage(context.Background(), params); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
	}

	return err
}

func (n *Notifier) generateProfileImage(req ImageRequest) ([]byte, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling image request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", n.imageServiceURL+"/api/image", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error sending request to image service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image service returned non-OK status: %d", resp.StatusCode)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading image response: %w", err)
	}

	return imageBytes, nil
}

func (n *Notifier) NotifyCollaborationVerified(collab db.Collaboration) error {
	var msgText string
	if collab.User.LanguageCode == db.LanguageRU {
		msgText = fmt.Sprintf("🎉 Поздравляем! Ваша коллаборация \"%s\" была подтверждена.\nТеперь она доступна всем пользователям Peatch.\n\nПосмотреть коллаборацию: %s/collaborations/%s", collab.Title, n.webappURL, collab.ID)
	} else {
		msgText = fmt.Sprintf("🎉 Congratulations! Your collaboration \"%s\" has been verified.\nIt is now visible to all Peatch users.\n\nView collaboration: %s/collaborations/%s", collab.Title, n.webappURL, collab.ID)
	}

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		ChatID: fmt.Sprintf("%d", collab.User.ChatID),
		Text:   msgText,
	})

	return err
}

func (n *Notifier) NotifyNewPendingUser(user db.User) error {
	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}

	msgText := fmt.Sprintf("🔔 New user pending verification:\nName: %s %s\nUsername: @%s\n\nVerify them in the admin panel: %s/admin/users",
		firstName, lastName, user.Username, n.adminWebApp)

	params := &telegram.SendMessageParams{ChatID: fmt.Sprintf("%d", n.adminChatID), Text: msgText}

	_, err := n.bot.SendMessage(context.Background(), params)

	return err
}

func (n *Notifier) NotifyNewPendingCollaboration(user db.User, collab db.Collaboration) error {
	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}

	msgText := fmt.Sprintf("🔔 New collaboration pending verification:\nTitle: %s\nBy: %s %s (@%s)\n\nVerify it in the admin panel: %s/admin/collaborations",
		collab.Title, firstName, lastName, user.Username, n.adminWebApp)

	params := &telegram.SendMessageParams{
		ChatID: fmt.Sprintf("%d", n.adminChatID),
		Text:   msgText,
	}

	_, err := n.bot.SendMessage(context.Background(), params)

	return err
}
