package job

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/notification"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type storage interface {
	GetUserProfile(params db.GetUsersParams) (*db.User, error)
	CreateNotification(notification db.Notification) (*db.Notification, error)
	SearchNotification(params db.NotificationQuery) (*db.Notification, error)
	ListUserCollaborations(from time.Time) ([]db.UserCollaborationRequest, error)
	UpdateNotificationSentAt(notificationID int64) error
	ListCollaborations(params db.CollaborationQuery) ([]db.Collaboration, error)
	FindMatchingUsers(exclude int64, opportunityIDs []int64, badgeIDs []int64) ([]db.User, error)
	ListNewUserProfiles(from time.Time) ([]db.User, error)
	ListCollaborationRequests(from time.Time) ([]db.CollaborationRequest, error)
	GetCollaborationOwner(collaborationID int64) (*db.User, error)
}

type notifyJob struct {
	storage  storage
	notifier notifier
	config   config
}

type notifier interface {
	SendNotification(params notification.SendNotificationParams) error
}

type config struct {
	imgServiceURL string
	botWebApp     string
	webappURL     string
	groupChatID   int64
}

func WithConfig(imgServiceURL, botWebApp, webappURL string, groupChatID int64) func(*notifyJob) {
	return func(j *notifyJob) {
		j.config = config{
			imgServiceURL: imgServiceURL,
			botWebApp:     botWebApp,
			webappURL:     webappURL,
			groupChatID:   groupChatID,
		}
	}
}

func NewNotifyJob(storage storage, notifier notifier, opts ...func(*notifyJob)) *notifyJob {
	j := &notifyJob{
		storage:  storage,
		notifier: notifier,
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

func (j *notifyJob) NotifyNewUserProfile() error {
	// Here fetch latest users that published their profile. Find matching users and send them a notification
	log.Println("Checking for new user profiles")

	newUsers, err := j.storage.ListNewUserProfiles(time.Now().Add(-24 * time.Hour))

	if err != nil {
		return err
	}

	for _, user := range newUsers {
		q := db.NotificationQuery{
			ChatID:           &j.config.groupChatID,
			NotificationType: db.NotificationTypeUserPublished,
			EntityType:       "users",
			EntityID:         user.ID,
		}

		_, err := j.storage.SearchNotification(q)

		if err != nil && errors.Is(err, db.ErrNotFound) {
			params := db.GetUsersParams{UserID: user.ID}

			userDetails, err := j.storage.GetUserProfile(params)

			if err != nil {
				return err
			}

			opportunityIDs := make([]int64, len(userDetails.Opportunities))

			for i, opportunity := range userDetails.Opportunities {
				opportunityIDs[i] = opportunity.ID
			}

			img, err := fetchPreviewImage(j.config.imgServiceURL, userDetails)

			if err != nil {
				return err
			}

			ntf := &db.Notification{
				NotificationType: db.NotificationTypeUserPublished,
				EntityType:       "users",
				EntityID:         user.ID,
				ChatID:           j.config.groupChatID,
			}

			text := fmt.Sprintf("Someone has just published a new profile")

			created, err := j.storage.CreateNotification(*ntf)

			if err != nil {
				return err
			}

			linkToProfile := fmt.Sprintf("%s?startapp=t-users-%s", j.config.botWebApp, user.Username)

			n := notification.SendNotificationParams{
				ChatID:    created.ChatID,
				Message:   text,
				BotWebApp: linkToProfile,
				Image:     img,
			}

			if err = j.notifier.SendNotification(n); err != nil {
				log.Printf("Failed to send notification to user %d", user.ID)
				return err
			}

			if err = j.storage.UpdateNotificationSentAt(created.ID); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (j *notifyJob) NotifyNewCollaboration() error {
	log.Println("Checking for new collaborations")

	dayAgo := time.Now().Add(-24 * time.Hour)

	newCollaborations, err := j.storage.ListCollaborations(db.CollaborationQuery{
		Limit:   10,
		Page:    1,
		From:    &dayAgo,
		Visible: true,
	})

	if err != nil {
		return err
	}

	for _, collaboration := range newCollaborations {
		q := db.NotificationQuery{
			ChatID:           &j.config.groupChatID,
			NotificationType: db.NotificationTypeCollaborationPublished,
			EntityType:       "collaborations",
			EntityID:         collaboration.ID,
		}

		_, err := j.storage.SearchNotification(q)

		if err != nil && errors.Is(err, db.ErrNotFound) {
			params := db.GetUsersParams{
				UserID: collaboration.UserID,
			}

			creator, err := j.storage.GetUserProfile(params)

			if err != nil {
				return err
			}

			img, err := fetchPreviewImage(j.config.imgServiceURL, creator)

			if err != nil {
				return err
			}

			text := fmt.Sprintf(
				"*%s has just posted a new opportunity* %s\n%s\n%s",
				bot.EscapeMarkdown(*creator.FirstName),
				bot.EscapeMarkdown(collaboration.Title),
				bot.EscapeMarkdown(collaboration.Description),
				bot.EscapeMarkdown(collaboration.GetLocation()),
			)

			ntf := &db.Notification{
				NotificationType: db.NotificationTypeCollaborationPublished,
				EntityType:       "collaborations",
				EntityID:         collaboration.ID,
				ChatID:           j.config.groupChatID,
			}

			created, err := j.storage.CreateNotification(*ntf)

			if err != nil {
				return err
			}

			linkToCollaboration := fmt.Sprintf("%s?startapp=t-collaborations-%d", j.config.botWebApp, collaboration.ID)

			n := notification.SendNotificationParams{
				ChatID:    created.ChatID,
				Message:   text,
				BotWebApp: linkToCollaboration,
				Image:     img,
			}

			if err = j.notifier.SendNotification(n); err != nil {
				log.Printf("Failed to send notification to user %d", collaboration.UserID)
				return err
			}

			if err = j.storage.UpdateNotificationSentAt(created.ID); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (j *notifyJob) NotifyMatchedCollaboration() error {
	log.Println("Checking for new collaborations, sending to matching users")

	dayAgo := time.Now().Add(-24 * time.Hour)

	newCollaborations, err := j.storage.ListCollaborations(db.CollaborationQuery{
		Limit:   10,
		Page:    1,
		From:    &dayAgo,
		Visible: true,
	})

	if err != nil {
		return err
	}

	for _, collaboration := range newCollaborations {
		receivers, err := j.storage.FindMatchingUsers(collaboration.UserID, []int64{collaboration.OpportunityID}, []int64{})

		if err != nil {
			return err
		}

		if len(receivers) == 0 {
			log.Printf("No users found for collaboration %d", collaboration.ID)
			continue
		}

		params := db.GetUsersParams{
			UserID: collaboration.UserID,
		}

		creator, err := j.storage.GetUserProfile(params)

		if err != nil {
			return err
		}

		img, err := fetchPreviewImage(j.config.imgServiceURL, creator)

		if err != nil {
			return err
		}

		for _, receiver := range receivers {
			q := db.NotificationQuery{
				UserID:           &receiver.ID,
				NotificationType: db.NotificationTypeCollaborationPublished,
				EntityType:       "collaborations",
				EntityID:         collaboration.ID,
			}

			_, err := j.storage.SearchNotification(q)
			if err != nil && errors.Is(err, db.ErrNotFound) {
				ntf := &db.Notification{
					UserID:           receiver.ID,
					NotificationType: db.NotificationTypeCollaborationPublished,
					EntityType:       "collaborations",
					EntityID:         collaboration.ID,
					ChatID:           receiver.ChatID,
				}

				text := fmt.Sprintf(
					"*%s has just posted a new opportunity* %s\n%s\n%s",
					bot.EscapeMarkdown(*creator.FirstName),
					bot.EscapeMarkdown(collaboration.Title),
					bot.EscapeMarkdown(collaboration.Description),
					bot.EscapeMarkdown(collaboration.GetLocation()),
				)

				created, err := j.storage.CreateNotification(*ntf)

				if err != nil {
					return err
				}

				linkToCollaboration := fmt.Sprintf("%s/collaborations/%d", j.config.webappURL, collaboration.ID)

				params := notification.SendNotificationParams{
					ChatID:    created.ChatID,
					Message:   text,
					WebAppURL: linkToCollaboration,
					Image:     img,
				}

				if err = j.notifier.SendNotification(params); err != nil {
					log.Printf("Failed to send notification to user %d", collaboration.UserID)
					return err
				}

				if err = j.storage.UpdateNotificationSentAt(created.ID); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		}
	}

	return nil
}

func (j *notifyJob) NotifyUserReceivedCollaborationRequest() error {
	log.Println("Checking for new user collaboration requests")

	newCollaborations, err := j.storage.ListUserCollaborations(time.Now().Add(-24 * time.Hour))

	if err != nil {
		return err
	}

	for _, collaboration := range newCollaborations {
		// check if  user already received exact same notification
		q := db.NotificationQuery{
			UserID:           &collaboration.UserID,
			NotificationType: db.NotificationTypeUserCollaboration,
			EntityType:       "user_collaboration_requests",
			EntityID:         collaboration.ID,
		}

		_, err := j.storage.SearchNotification(q)

		if err != nil && errors.Is(err, db.ErrNotFound) {

			requester, err := j.storage.GetUserProfile(db.GetUsersParams{UserID: collaboration.RequesterID})

			if err != nil {
				return err
			}

			receiver, err := j.storage.GetUserProfile(db.GetUsersParams{UserID: collaboration.UserID})

			if err != nil {
				return err
			}

			img, err := fetchPreviewImage(j.config.imgServiceURL, requester)

			if err != nil {
				return err
			}

			ntf := &db.Notification{
				UserID:           collaboration.UserID,
				NotificationType: db.NotificationTypeUserCollaboration,
				EntityType:       "user_collaboration_requests",
				EntityID:         collaboration.ID,
				ChatID:           receiver.ChatID,
			}

			text := fmt.Sprintf(
				"%s sends you a message:\n%s",
				bot.EscapeMarkdown(*requester.FirstName),
				bot.EscapeMarkdown(collaboration.Message),
			)

			created, err := j.storage.CreateNotification(*ntf)

			if err != nil {
				return err
			}

			linkToProfile := fmt.Sprintf("%s/users/%d", j.config.webappURL, requester.Username)

			params := notification.SendNotificationParams{
				ChatID:    created.ChatID,
				Message:   text,
				WebAppURL: linkToProfile,
				Image:     img,
				Username:  &requester.Username,
			}

			if err = j.notifier.SendNotification(params); err != nil {
				log.Printf("Failed to send notification to user %d", collaboration.UserID)
				return err
			}

			if err = j.storage.UpdateNotificationSentAt(created.ID); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (j *notifyJob) NotifyCollaborationRequest() error {
	// Here fetch latest collaboration requests. Send notification to user who received the request on their collaboration
	log.Println("Checking for new collaboration requests")

	newRequests, err := j.storage.ListCollaborationRequests(time.Now().Add(-24 * time.Hour))

	if err != nil {
		return err
	}

	for _, request := range newRequests {
		// get the one who created the collaboration
		creator, err := j.storage.GetCollaborationOwner(request.CollaborationID)

		if err != nil {
			return err
		}

		// check if  user already received exact same notification
		q := db.NotificationQuery{
			UserID:           &creator.ID,
			NotificationType: db.NotificationTypeCollaborationRequest,
			EntityType:       "collaboration_requests",
			EntityID:         request.ID,
		}

		if _, err := j.storage.SearchNotification(q); err != nil && errors.Is(err, db.ErrNotFound) {
			// get the one who created the request
			requester, err := j.storage.GetUserProfile(db.GetUsersParams{UserID: request.UserID})

			if err != nil {
				return err
			}

			img, err := fetchPreviewImage(j.config.imgServiceURL, requester)

			if err != nil {
				return err
			}

			ntf := &db.Notification{
				UserID:           creator.ID,
				NotificationType: db.NotificationTypeCollaborationRequest,
				EntityType:       "collaboration_requests",
				EntityID:         request.ID,
				ChatID:           creator.ChatID,
			}

			text := fmt.Sprintf(
				"*%s wants to collaborate with you on your opportunity*\n%s",
				bot.EscapeMarkdown(*requester.FirstName),
				bot.EscapeMarkdown(request.Message))

			created, err := j.storage.CreateNotification(*ntf)

			if err != nil {
				return err
			}

			linkToProfile := fmt.Sprintf("%s/users/%d", j.config.webappURL, requester.Username)

			params := notification.SendNotificationParams{
				ChatID:    creator.ChatID,
				Message:   text,
				WebAppURL: linkToProfile,
				Image:     img,
				Username:  &requester.Username,
			}

			if err = j.notifier.SendNotification(params); err != nil {
				log.Printf("Failed to send notification to user %d", creator.ID)
			}

			if err = j.storage.UpdateNotificationSentAt(created.ID); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
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

func fetchPreviewImage(baseUrl string, user *db.User) ([]byte, error) {
	if user.AvatarURL == nil || user.FirstName == nil || user.LastName == nil || user.Title == nil {
		return nil, errors.New("user data is not complete")
	}

	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	u.Path += "/api/image"

	body := ImageRequest{
		Title:    fmt.Sprintf("%s %s", *user.FirstName, *user.LastName),
		Subtitle: *user.Title,
		Avatar:   fmt.Sprintf("https://assets.peatch.io/%s", *user.AvatarURL),
	}

	tags := make([]Tag, len(user.Badges))

	for i, badge := range user.Badges {
		tags[i] = Tag{
			Text:  badge.Text,
			Color: badge.Color,
			Icon:  badge.Icon,
		}
	}

	body.Tags = tags

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to download image: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download image, got status code: %d", resp.StatusCode)
		return nil, errors.New(fmt.Sprintf("failed to download image, status code: %d", resp.StatusCode))
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read image data: %s", err)
		return nil, err
	}

	return imgData, nil
}
