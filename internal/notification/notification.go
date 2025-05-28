package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/peatch-io/peatch/internal/db"
)

type NotifierConfig struct {
	BotToken         string
	AdminChatID      int64
	CommunityChatID  int64
	BotWebApp        string
	WebAppURL        string
	AdminWebApp      string
	ImageServiceURL  string
	TestNotification bool
}

type Notifier struct {
	bot              *telegram.Bot
	adminChatID      int64
	communityChatID  int64
	botWebApp        string
	webappURL        string
	adminWebApp      string
	imageServiceURL  string
	testNotification bool
}

func NewNotifier(config NotifierConfig, bot *telegram.Bot) *Notifier {
	return &Notifier{
		bot:              bot,
		adminChatID:      config.AdminChatID,
		communityChatID:  config.CommunityChatID,
		botWebApp:        config.BotWebApp,
		webappURL:        config.WebAppURL,
		adminWebApp:      config.AdminWebApp,
		imageServiceURL:  config.ImageServiceURL,
		testNotification: config.TestNotification,
	}
}

// getChatID returns the appropriate chat ID based on test mode
func (n *Notifier) getChatID(originalChatID interface{}) string {
	if n.testNotification {
		return fmt.Sprintf("%d", n.adminChatID)
	}

	switch v := originalChatID.(type) {
	case int64:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

type ImageRequest struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Avatar   string `json:"avatar"`
	Tags     []Tag  `json:"tags"`
}

type CollaborationImageRequest struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Tags     []Tag    `json:"tags"`
	User     UserInfo `json:"user"`
}

type UserInfo struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}

type Tag struct {
	Text  string `json:"text"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func (n *Notifier) NotifyUserVerified(user db.User) error {
	var msgText string
	if user.LanguageCode == db.LanguageRU {
		msgText = fmt.Sprintf("🎉 Поздравляем\\! Ваш профиль был подтверждён\\. 💡 Следующие шаги:\n• Ищешь кого\\-то? Запости \\- и мы сообщим подходящим людям\\.\n• [Вступай в комьюнити](https://t.me/peatch_community), чтобы быть в курсе событий\\.")
	} else {
		msgText = fmt.Sprintf("🎉 Congratulations\\! Your profile has been verified\\.\n• Looking for someone? Post it \\- we'll notify the right people\\.\n• [Join the community](https://t.me/peatch_community) to stay updated\\.")
	}

	btnText := "Publish Collaboration"
	if user.LanguageCode == db.LanguageRU {
		btnText = "Опубликовать проект"
	}

	button := models.InlineKeyboardButton{
		Text: btnText,
		WebApp: &models.WebAppInfo{
			URL: fmt.Sprintf("%s/collaborations/edit", n.webappURL),
		},
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		ChatID:      n.getChatID(user.ChatID),
		Text:        msgText,
		ReplyMarkup: &keyboard,
		ParseMode:   models.ParseModeMarkdown,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: new(bool), // Disable link previews
		},
	})

	if err != nil {
		return err
	}

	if !user.IsProfileComplete() {
		log.Printf("User %s has an incomplete profile. Not sending welcome message.", user.Username)
		return nil
	}

	fullName := ""
	if user.Name != nil {
		fullName = *user.Name
	} else {
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
		imageReq.Avatar = fmt.Sprintf("https://assets.peatch.io/cdn-cgi/image/width=400/%s", *user.AvatarURL)
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

	btnText = "View Profile"

	button = models.InlineKeyboardButton{
		Text: btnText,
		URL:  fmt.Sprintf("%s?startapp=u_%s", n.botWebApp, user.ID),
	}

	keyboard = models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	if imageBytes != nil && len(imageBytes) > 0 {
		photoData := &models.InputFileUpload{
			Filename: fmt.Sprintf("welcome_%s.png", user.ID),
			Data:     bytes.NewReader(imageBytes),
		}

		if _, err := n.bot.SendPhoto(context.Background(), &telegram.SendPhotoParams{
			ChatID:      n.getChatID(n.communityChatID),
			Caption:     communityMsg,
			Photo:       photoData,
			ReplyMarkup: &keyboard,
		}); err != nil {
			fmt.Printf("Error sending photo: %v\n", err)
		}
	} else {
		params := &telegram.SendMessageParams{
			ChatID:      n.getChatID(n.communityChatID),
			Text:        communityMsg,
			ReplyMarkup: &keyboard,
		}

		if _, err := n.bot.SendMessage(context.Background(), params); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
	}

	return err
}

func (n *Notifier) generateImage(endpoint string, request interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", n.imageServiceURL+endpoint, bytes.NewBuffer(jsonData))
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

func (n *Notifier) generateProfileImage(req ImageRequest) ([]byte, error) {
	return n.generateImage("/api/image", req)
}

func (n *Notifier) generateCollaborationImage(req CollaborationImageRequest) ([]byte, error) {
	return n.generateImage("/api/collaboration", req)
}

func (n *Notifier) NotifyCollaborationVerified(collab db.Collaboration) error {
	var msgText string
	if collab.User.LanguageCode == db.LanguageRU {
		msgText = fmt.Sprintf("🎉 Поздравляем! Ваша коллаборация \"%s\" была подтверждена.\nТеперь она доступна всем пользователям Peatch.", collab.Title)
	} else {
		msgText = fmt.Sprintf("🎉 Congratulations! Your collaboration \"%s\" has been verified.\nIt is now visible to all Peatch users.", collab.Title)
	}

	btnText := "View Collaboration"
	if collab.User.LanguageCode == db.LanguageRU {
		btnText = "Посмотреть коллаборацию"
	}

	button := models.InlineKeyboardButton{
		Text: btnText,
		URL:  fmt.Sprintf("%s?startapp=c_%s", n.botWebApp, collab.ID),
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		ChatID:      n.getChatID(collab.User.ChatID),
		Text:        msgText,
		ReplyMarkup: &keyboard,
	})

	if err != nil {
		return err
	}

	return n.SendCollaborationToCommunityChatWithImage(collab)
}

// NotifyUsersWithMatchingOpportunity sends a notification to all users with a matching opportunity
// Uses batching and rate limiting to respect Telegram's limit of 30 messages per second
func (n *Notifier) NotifyUsersWithMatchingOpportunity(collab db.Collaboration, users []db.User) error {
	if len(users) == 0 {
		return nil
	}

	// Create the notification button
	button := models.InlineKeyboardButton{
		Text: "View Collaboration",
		URL:  fmt.Sprintf("%s?startapp=c_%s", n.botWebApp, collab.ID),
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	// Prepare a batch of users to notify
	eligibleUsers := make([]db.User, 0, len(users))
	for _, user := range users {
		// Skip the collaboration owner
		if user.ID == collab.UserID {
			continue
		}

		// Skip users with notifications disabled
		//if user.NotificationsEnabledAt == nil {
		//	continue
		//}

		eligibleUsers = append(eligibleUsers, user)
	}

	fmt.Printf("Found %d eligible users for collaboration %s\n", len(eligibleUsers), collab.ID)

	// If no eligible users, return early
	if len(eligibleUsers) == 0 {
		return nil
	}

	// Generate collaboration image once for all notifications if image service is available
	var imageBytes []byte
	if n.imageServiceURL != "" {
		// Prepare image generation request
		fullName := ""
		if collab.User.Name != nil {
			fullName = *collab.User.Name
		} else {
			fullName = collab.User.Username
		}

		// Extract badges for tags
		tags := make([]Tag, 0, 5)
		if collab.Badges != nil && len(collab.Badges) > 0 {
			for i, badge := range collab.Badges {
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

		userRole := "Member"
		if collab.User.Title != nil {
			userRole = *collab.User.Title
		}

		userAvatarURL := ""
		if collab.User.AvatarURL != nil {
			userAvatarURL = fmt.Sprintf("https://assets.peatch.io/cdn-cgi/image/width=400/%s", *collab.User.AvatarURL)
		}

		// Create the image request
		imageReq := CollaborationImageRequest{
			Title:    collab.Title,
			Subtitle: collab.Opportunity.Text,
			Tags:     tags,
			User: UserInfo{
				Avatar: userAvatarURL,
				Name:   fullName,
				Role:   userRole,
			},
		}

		// Generate the image
		var err error
		imageBytes, err = n.generateCollaborationImage(imageReq)
		if err != nil {
			log.Printf("Error generating collaboration image: %v", err)
			// Continue with text-only notifications if image generation fails
		}
	}

	// Use a channel to control rate limiting
	// Telegram API limit is approximately 30 messages per second
	const telegramRateLimit = 30
	const batchSize = 15
	const cooldownPeriod = time.Second // 1 second between batches

	// Process users in batches
	for i := 0; i < len(eligibleUsers); i += batchSize {
		end := i + batchSize
		if end > len(eligibleUsers) {
			end = len(eligibleUsers)
		}

		batch := eligibleUsers[i:end]

		// Start a goroutine for the entire batch to process in parallel
		go func(batchUsers []db.User, imageData []byte) {
			rateLimiter := time.NewTicker(time.Second / telegramRateLimit)
			defer rateLimiter.Stop()

			for _, user := range batchUsers {
				<-rateLimiter.C // Wait for rate limiter before sending next message

				var msgText string
				if user.LanguageCode == db.LanguageRU {
					collabUserName := collab.User.Username
					if collab.User.Name != nil {
						collabUserName = *collab.User.Name
					}
					msgText = fmt.Sprintf("🔍 Новая коллаборация от %s, которая может вас заинтересовать!\n\n%s",
						collabUserName,
						collab.Title)
				} else {
					collabUserName := collab.User.Username
					if collab.User.Name != nil {
						collabUserName = *collab.User.Name
					}
					msgText = fmt.Sprintf("🔍 New collaboration from %s that might interest you!\n\n%s",
						collabUserName,
						collab.Title)
				}

				var err error
				// If we have image bytes, send as photo with caption, otherwise send as text message
				if imageData != nil && len(imageData) > 0 {
					photoData := &models.InputFileUpload{
						Filename: fmt.Sprintf("collab_opportunity_%s.png", collab.ID),
						Data:     bytes.NewReader(imageData),
					}

					_, err = n.bot.SendPhoto(context.Background(), &telegram.SendPhotoParams{
						ChatID:      n.getChatID(user.ChatID),
						Caption:     msgText,
						Photo:       photoData,
						ReplyMarkup: &keyboard,
					})
					fmt.Printf("Sent opportunity match notification to user %s with image\n", user.ID)
				} else {
					// Fall back to text-only message if no image data available
					_, err = n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
						ChatID:      n.getChatID(user.ChatID),
						Text:        msgText,
						ReplyMarkup: &keyboard,
					})
					fmt.Printf("Sent opportunity match notification to user %s without image\n", user.ID)
				}

				if err != nil {
					log.Printf("Failed to send opportunity match notification to user %s: %v", user.ID, err)
				}
			}
		}(batch, imageBytes)

		// Add a cooldown between batches to ensure we don't exceed global rate limits
		if end < len(eligibleUsers) {
			time.Sleep(cooldownPeriod)
		}
	}

	return nil
}

func (n *Notifier) SendCollaborationToCommunityChatWithImage(collab db.Collaboration) error {
	fullName := ""
	if collab.User.Name != nil {
		fullName = *collab.User.Name
	} else {
		fullName = collab.User.Username
	}

	button := models.InlineKeyboardButton{
		Text: "View Collaboration",
		URL:  fmt.Sprintf("%s?startapp=c_%s", n.botWebApp, collab.ID),
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	tags := make([]Tag, 0, 5)
	if collab.Badges != nil && len(collab.Badges) > 0 {
		for i, badge := range collab.Badges {
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

	userRole := "Member"
	if collab.User.Title != nil {
		userRole = *collab.User.Title
	}

	userAvatarURL := ""
	if collab.User.AvatarURL != nil {
		userAvatarURL = fmt.Sprintf("https://assets.peatch.io/cdn-cgi/image/width=400/%s", *collab.User.AvatarURL)
	}

	imageReq := CollaborationImageRequest{
		Title:    collab.Title,
		Subtitle: collab.Opportunity.Text,
		Tags:     tags,
		User: UserInfo{
			Avatar: userAvatarURL,
			Name:   fullName,
			Role:   userRole,
		},
	}

	var imageBytes []byte
	var err error
	if n.imageServiceURL != "" {
		imageBytes, err = n.generateCollaborationImage(imageReq)
		if err != nil {
			fmt.Printf("Error generating collaboration image: %v\n", err)
		}
	}

	communityMsg := fmt.Sprintf("🌟 New collaboration opportunity!\n\"%s\" by %s\n\nCheck it out",
		collab.Title, fullName)

	if imageBytes != nil && len(imageBytes) > 0 {
		photoData := &models.InputFileUpload{
			Filename: fmt.Sprintf("collab_%s.png", collab.ID),
			Data:     bytes.NewReader(imageBytes),
		}

		if _, err := n.bot.SendPhoto(context.Background(), &telegram.SendPhotoParams{
			ChatID:      n.getChatID(n.communityChatID),
			Caption:     communityMsg,
			Photo:       photoData,
			ReplyMarkup: &keyboard,
		}); err != nil {
			fmt.Printf("Error sending photo: %v\n", err)
			return err
		}
	} else {
		params := &telegram.SendMessageParams{
			ChatID:      n.getChatID(n.communityChatID),
			Text:        communityMsg,
			ReplyMarkup: &keyboard,
		}

		if _, err := n.bot.SendMessage(context.Background(), params); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			return err
		}
	}

	return nil
}

func (n *Notifier) NotifyNewPendingUser(user db.User) error {
	name := ""
	if user.Name != nil {
		name = *user.Name
	} else {
		name = user.Username
	}

	btn := models.InlineKeyboardButton{
		Text: "View Profile",
		URL:  n.adminWebApp,
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{btn},
		},
	}

	msgText := fmt.Sprintf("🔔 New user pending verification:\nName: %s\nUsername: @%s",
		name, user.Username)

	params := &telegram.SendMessageParams{
		ChatID:      n.getChatID(n.adminChatID),
		Text:        msgText,
		ReplyMarkup: &keyboard,
	}

	_, err := n.bot.SendMessage(context.Background(), params)

	return err
}

func (n *Notifier) NotifyNewPendingCollaboration(collab db.Collaboration) error {
	user := collab.User

	name := ""
	if user.Name != nil {
		name = *user.Name
	}

	btn := models.InlineKeyboardButton{
		Text: "View Collaboration",
		URL:  n.adminWebApp,
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{btn},
		},
	}

	msgText := fmt.Sprintf("🔔 New collaboration pending verification:\nTitle: %s\nBy: %s (@%s)",
		collab.Title, name, user.Username)

	params := &telegram.SendMessageParams{
		ChatID:      n.getChatID(n.adminChatID),
		Text:        msgText,
		ReplyMarkup: &keyboard,
	}

	_, err := n.bot.SendMessage(context.Background(), params)

	return err
}

var ErrUserBlockedBot = errors.New("user has blocked the bot")

func (n *Notifier) NotifyUserFollow(userToFollow db.User, follower db.User) error {

	if userToFollow.ChatID == 0 {
		return fmt.Errorf("user to follow %s has no chat ID", userToFollow.ID)
	}

	followerName := follower.Username
	if follower.Name != nil {
		followerName = *follower.Name
	}

	var msgText string
	if userToFollow.LanguageCode == db.LanguageRU {
		msgText = fmt.Sprintf("👋 [%s](https://t.me/%s) заметил ваш профиль, напишите ему в Telegram\\!", telegram.EscapeMarkdown(followerName), follower.Username)
	} else {
		msgText = fmt.Sprintf("👋 [%s](https://t.me/%s) noticed your profile, write to him in Telegram\\!", telegram.EscapeMarkdown(followerName), follower.Username)
	}

	btnText := "View Profile"
	if userToFollow.LanguageCode == db.LanguageRU {
		btnText = "Посмотреть профиль"
	}

	button := models.InlineKeyboardButton{
		Text: btnText,
		URL:  fmt.Sprintf("%s?startapp=u_%s", n.botWebApp, follower.ID),
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	disabled := true

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		ChatID: n.getChatID(userToFollow.ChatID),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &disabled,
		},
		Text:        msgText,
		ReplyMarkup: &keyboard,
		ParseMode:   models.ParseModeMarkdown,
	})

	if err != nil && strings.Contains(err.Error(), "Forbidden") &&
		(strings.Contains(err.Error(), "bot was blocked by the user") ||
			strings.Contains(err.Error(), "user is deactivated")) {
		return ErrUserBlockedBot
	}

	return err
}

func (n *Notifier) NotifyUserVerificationDenied(user db.User) error {
	var msgText string
	if user.LanguageCode == db.LanguageRU {
		msgText = fmt.Sprintf("⚠️ Твой профиль не прошел проверку.\n\nПохоже, он слишком пустой или похож на спам. Добавь больше деталей и искренности, и попробуй снова!")
	} else {
		msgText = fmt.Sprintf("⚠️ Your profile didn’t pass verification.\nIt seems too empty or spammy. Add more details and sincerity, and try again!")
	}

	btnText := "Update Profile"
	if user.LanguageCode == db.LanguageRU {
		btnText = "Обновить профиль"
	}

	button := models.InlineKeyboardButton{
		Text:   btnText,
		WebApp: &models.WebAppInfo{URL: fmt.Sprintf("%s/users/edit", n.webappURL)},
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		ChatID:      n.getChatID(user.ChatID),
		Text:        msgText,
		ReplyMarkup: &keyboard,
	})

	return err
}

func (n *Notifier) NotifyCollaborationVerificationDenied(collab db.Collaboration) error {
	var msgText string
	if collab.User.LanguageCode == db.LanguageRU {
		msgText = fmt.Sprintf("⚠️ Ваша коллаборация \"%s\" не прошла проверку.\nПохоже, описание слишком общее или неубедительное. Добавьте больше деталей и конкретики, и попробуйте снова!", collab.Title)
	} else {
		msgText = fmt.Sprintf("⚠️ Your collaboration \"%s\" didn’t pass verification.\nThe description seems too vague or unconvincing. Add more details and specifics, and try again!", collab.Title)
	}

	btnText := "Update Collaboration"
	if collab.User.LanguageCode == db.LanguageRU {
		btnText = "Обновить коллаборацию"
	}

	button := models.InlineKeyboardButton{
		Text:   btnText,
		WebApp: &models.WebAppInfo{URL: fmt.Sprintf("%s/collaborations/edit", n.webappURL)},
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		ChatID:      n.getChatID(collab.User.ChatID),
		Text:        msgText,
		ReplyMarkup: &keyboard,
	})

	return err
}

func (n *Notifier) NotifyCollabInterest(collab db.Collaboration, user db.User) error {
	if collab.User.ChatID == 0 {
		return fmt.Errorf("collaboration owner %s has no chat ID", collab.User.ID)
	}

	userName := user.Username
	if user.Name != nil && *user.Name != "" {
		userName = *user.Name
	}

	var msgText string
	if collab.User.LanguageCode == db.LanguageRU {
		msgText = fmt.Sprintf("🔔 [%s](https://t.me/%s) проявил интерес к вашему проекту \"%s\", напишите ему в Telegram\\!", telegram.EscapeMarkdown(userName), user.Username, collab.Title)
	} else {
		msgText = fmt.Sprintf("🔔 [%s](https://t.me/%s) expressed interest in your project \"%s\", write to them in Telegram\\!", telegram.EscapeMarkdown(userName), user.Username, collab.Title)
	}

	btnText := "View Profile"
	if collab.User.LanguageCode == db.LanguageRU {
		btnText = "Посмотреть профиль"
	}

	button := models.InlineKeyboardButton{
		Text: btnText,
		URL:  fmt.Sprintf("%s?startapp=u_%s", n.botWebApp, user.ID),
	}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{button},
		},
	}

	disabled := true

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		ChatID: n.getChatID(collab.User.ChatID),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &disabled,
		},
		Text:        msgText,
		ReplyMarkup: &keyboard,
		ParseMode:   models.ParseModeMarkdown,
	})

	if err != nil && strings.Contains(err.Error(), "Forbidden") &&
		(strings.Contains(err.Error(), "bot was blocked by the user") ||
			strings.Contains(err.Error(), "user is deactivated")) {
		return ErrUserBlockedBot
	}

	return err
}
