package job

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/peatch-io/peatch/internal/db"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type storage interface {
	GetUserByID(id int64, showHidden bool) (*db.User, error)
	CreateNotification(notification db.Notification) (*db.Notification, error)
	SearchNotification(params db.NotificationQuery) (*db.Notification, error)
	ListUserCollaborations(from time.Time) ([]db.UserCollaborationRequest, error)
	UpdateNotificationSentAt(notificationID int64) error
	ListCollaborations(params db.CollaborationQuery) ([]db.Collaboration, error)
	FindMatchingUsers(opportunityIDs []int64, badgeIDs []int64) ([]db.User, error)
	ListNewUserProfiles(from time.Time) ([]db.User, error)
	ListCollaborationRequests(from time.Time) ([]db.CollaborationRequest, error)
	GetCollaborationOwner(collaborationID int64) (*db.User, error)
}

type notifyJob struct {
	storage       storage
	notifier      notifier
	imgServiceURL string
	webAppURL     string
	groupChatID   int64
}

type notifier interface {
	SendNotification(chatID int64, message string, link string, img []byte) error
}

func NewNotifyJob(storage storage, notifier notifier, imgServiceURL, webAppURL string, groupChatID int64) *notifyJob {
	return &notifyJob{
		storage:       storage,
		notifier:      notifier,
		imgServiceURL: imgServiceURL,
		webAppURL:     webAppURL,
		groupChatID:   groupChatID,
	}
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
			ChatID:           &j.groupChatID,
			NotificationType: db.NotificationTypeUserPublished,
			EntityType:       "users",
			EntityID:         user.ID,
		}

		_, err := j.storage.SearchNotification(q)

		if err != nil && errors.Is(err, db.ErrNotFound) {
			userDetails, err := j.storage.GetUserByID(user.ID, false)

			if err != nil {
				return err
			}

			opportunityIDs := make([]int64, len(userDetails.Opportunities))

			for i, opportunity := range userDetails.Opportunities {
				opportunityIDs[i] = opportunity.ID
			}

			img, err := fetchPreviewImage(j.imgServiceURL, userDetails)

			if err != nil {
				return err
			}

			notification := &db.Notification{
				NotificationType: db.NotificationTypeUserPublished,
				EntityType:       "users",
				EntityID:         user.ID,
				ChatID:           j.groupChatID,
			}

			text := fmt.Sprintf("Someone has just published a new profile")

			created, err := j.storage.CreateNotification(*notification)

			if err != nil {
				return err
			}

			linkToProfile := fmt.Sprintf("%s/users/%d", j.webAppURL, user.ID)

			if err = j.notifier.SendNotification(created.ChatID, text, linkToProfile, img); err != nil {
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
			ChatID:           &j.groupChatID,
			NotificationType: db.NotificationTypeCollaborationPublished,
			EntityType:       "collaborations",
			EntityID:         collaboration.ID,
		}

		_, err := j.storage.SearchNotification(q)

		if err != nil && errors.Is(err, db.ErrNotFound) {
			creator, err := j.storage.GetUserByID(collaboration.UserID, false)

			if err != nil {
				return err
			}

			img, err := fetchPreviewImage(j.imgServiceURL, creator)

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

			notification := &db.Notification{
				NotificationType: db.NotificationTypeCollaborationPublished,
				EntityType:       "collaborations",
				EntityID:         collaboration.ID,
				ChatID:           j.groupChatID,
			}

			created, err := j.storage.CreateNotification(*notification)

			if err != nil {
				return err
			}

			linkToCollaboration := fmt.Sprintf("%s/collaborations/%d", j.webAppURL, collaboration.ID)

			if err = j.notifier.SendNotification(created.ChatID, text, linkToCollaboration, img); err != nil {
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
			requester, err := j.storage.GetUserByID(collaboration.RequesterID, false)

			if err != nil {
				return err
			}

			receiver, err := j.storage.GetUserByID(collaboration.UserID, false)

			if err != nil {
				return err
			}

			img, err := fetchPreviewImage(j.imgServiceURL, requester)

			if err != nil {
				return err
			}

			notification := &db.Notification{
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

			created, err := j.storage.CreateNotification(*notification)

			if err != nil {
				return err
			}

			linkToProfile := fmt.Sprintf("%s/users/%d", j.webAppURL, collaboration.RequesterID)

			if err = j.notifier.SendNotification(created.ChatID, text, linkToProfile, img); err != nil {
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
			requester, err := j.storage.GetUserByID(request.UserID, false)

			if err != nil {
				return err
			}

			img, err := fetchPreviewImage(j.imgServiceURL, requester)

			if err != nil {
				return err
			}

			notification := &db.Notification{
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

			created, err := j.storage.CreateNotification(*notification)

			if err != nil {
				return err
			}

			linkToProfile := fmt.Sprintf("%s/users/%d", j.webAppURL, requester.ID)

			if err = j.notifier.SendNotification(created.ChatID, text, linkToProfile, img); err != nil {
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
