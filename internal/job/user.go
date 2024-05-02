package job

import (
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
	SearchNotification(userID int64, notificationType db.NotificationType, entityType string, entityID int64) (*db.Notification, error)
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
}

type notifier interface {
	SendNotification(chatID int64, message string, link string, img []byte) error
}

func NewNotifyJob(storage storage, notifier notifier, imgServiceURL, webAppURL string) *notifyJob {
	return &notifyJob{
		storage:       storage,
		notifier:      notifier,
		imgServiceURL: imgServiceURL,
		webAppURL:     webAppURL,
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
		userDetails, err := j.storage.GetUserByID(user.ID, true)

		if err != nil {
			return err
		}

		opportunityIDs := make([]int64, len(userDetails.Opportunities))

		for i, opportunity := range userDetails.Opportunities {
			opportunityIDs[i] = opportunity.ID
		}

		receivers, err := j.storage.FindMatchingUsers(opportunityIDs, []int64{})

		if err != nil {
			return err
		}

		if len(receivers) == 0 {
			log.Printf("No users that match user %d opportunities", user.ID)
			continue
		}

		img, err := fetchPreviewImage(j.imgServiceURL, userDetails)

		if err != nil {
			return err
		}

		for _, receiver := range receivers {
			_, err := j.storage.SearchNotification(
				receiver.ID,
				db.NotificationTypeUserPublished,
				"users",
				user.ID,
			)

			if err != nil && errors.Is(err, db.ErrNotFound) {

				notification := &db.Notification{
					UserID:           receiver.ID,
					NotificationType: db.NotificationTypeUserPublished,
					EntityType:       "users",
					EntityID:         user.ID,
					ChatID:           receiver.ChatID,
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
		receivers, err := j.storage.FindMatchingUsers([]int64{collaboration.OpportunityID}, []int64{})

		if err != nil {
			return err
		}

		if len(receivers) == 0 {
			log.Printf("No users found for collaboration %d", collaboration.ID)
			continue
		}

		creator, err := j.storage.GetUserByID(collaboration.UserID, true)

		if err != nil {
			return err
		}

		img, err := fetchPreviewImage(j.imgServiceURL, creator)

		if err != nil {
			return err
		}

		for _, receiver := range receivers {
			_, err := j.storage.SearchNotification(
				receiver.ID,
				db.NotificationTypeCollaborationPublished,
				"collaborations",
				collaboration.ID,
			)

			if err != nil && errors.Is(err, db.ErrNotFound) {

				notification := &db.Notification{
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
		_, err := j.storage.SearchNotification(
			collaboration.UserID,
			db.NotificationTypeUserCollaboration,
			"user_collaboration_requests",
			collaboration.ID,
		)

		if err != nil && errors.Is(err, db.ErrNotFound) {
			requester, err := j.storage.GetUserByID(collaboration.RequesterID, true)

			if err != nil {
				return err
			}

			receiver, err := j.storage.GetUserByID(collaboration.UserID, true)

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
		if _, err := j.storage.SearchNotification(
			creator.ID,
			db.NotificationTypeCollaborationRequest,
			"collaboration_requests",
			request.ID,
		); err != nil && errors.Is(err, db.ErrNotFound) {
			// get the one who created the request
			requester, err := j.storage.GetUserByID(request.UserID, true)

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

func fetchPreviewImage(baseUrl string, user *db.User) ([]byte, error) {
	if user.AvatarURL == nil || user.FirstName == nil || user.LastName == nil || user.Title == nil {
		return nil, errors.New("user data is not complete")
	}

	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	u.Path += "/api/image"

	params := url.Values{}
	params.Add("title", fmt.Sprintf("%s %s", *user.FirstName, *user.LastName))
	params.Add("subtitle", *user.Title)
	params.Add("avatar", fmt.Sprintf("https://assets.peatch.io/%s", *user.AvatarURL))

	if len(user.Badges) > 0 {
		tags := ""
		for _, badge := range user.Badges {
			tags += fmt.Sprintf("%s,%s,%s;", badge.Text, badge.Color, badge.Icon)
		}
		params.Add("tags", tags)
	}

	u.RawQuery = params.Encode()

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Get(u.String())
	if err != nil {
		log.Printf("Failed to download image: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download image, got status code: %d", resp.StatusCode)
		return nil, errors.New("failed to download image")
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read image data: %s", err)
		return nil, err
	}

	return imgData, nil
}
