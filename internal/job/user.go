package job

import (
	"errors"
	"fmt"
	"github.com/peatch-io/peatch/internal/db"
	"log"
	"time"
)

type storage interface {
	GetUserByID(id int64) (*db.User, error)
	CreateNotification(notification db.Notification) (*db.Notification, error)
	SearchNotification(userID int64, notificationType db.NotificationType, entityType string, entityID int64) (*db.Notification, error)
	ListCollaborations(query db.CollaborationQuery) ([]db.Collaboration, error)
	ListUserCollaborations(from time.Time) ([]db.UserCollaborationRequest, error)
	UpdateNotificationSentAt(notificationID int64) error
	ListNewCollaborations(from time.Time) ([]db.Collaboration, error)
	FindMatchingUsers(opportunityIDs []int64, badgeIDs []int64) ([]db.User, error)
}

type notifyJob struct {
	storage  storage
	notifier notifier
}

type notifier interface {
	SendNotification(chatID int64, message string, imgUrl string, link string) error
}

func NewNotifyJob(storage storage, notifier notifier) *notifyJob {
	return &notifyJob{
		storage:  storage,
		notifier: notifier,
	}
}

func (j *notifyJob) NotifyNewCollaboration() error {
	log.Println("Checking for new collaborations")

	newCollaborations, err := j.storage.ListNewCollaborations(time.Now().Add(-24 * time.Hour))

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

		creator, err := j.storage.GetUserByID(collaboration.UserID)

		if err != nil {
			return err
		}

		imgURL, err := createImgURL(creator)

		for _, receiver := range receivers {
			_, err := j.storage.SearchNotification(
				receiver.ID,
				db.NotificationTypeCollaborationPublished,
				"collaborations",
				collaboration.ID,
			)

			if err != nil && errors.As(err, &db.ErrNotFound) {

				notification := &db.Notification{
					UserID:           receiver.ID,
					NotificationType: db.NotificationTypeCollaborationPublished,
					Text:             fmt.Sprintf("%s wants to collaborate with you", *creator.FirstName),
					EntityType:       "collaborations",
					EntityID:         collaboration.ID,
					ChatID:           receiver.ChatID,
					ImageURL:         imgURL,
				}

				created, err := j.storage.CreateNotification(*notification)

				if err != nil {
					return err
				}

				linkToCollaboration := fmt.Sprintf("https://peatch.pages.dev/collaborations/%d", collaboration.ID)

				if err = j.notifier.SendNotification(created.ChatID, created.Text, imgURL, linkToCollaboration); err != nil {
					log.Printf("Failed to send notification to user %d", collaboration.UserID)
					return err
				}

				if err = j.storage.UpdateNotificationSentAt(created.ID); err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				log.Printf("Notification already sent for user %d and collaboration %d", collaboration.UserID, collaboration.ID)
			}
		}
	}

	return nil
}

func (j *notifyJob) NotifyUserReceivedCollaborationRequest() error {
	log.Println("Checking for new collaboration requests")

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

		if err != nil && errors.As(err, &db.ErrNotFound) {
			requester, err := j.storage.GetUserByID(collaboration.RequesterID)

			if err != nil {
				return err
			}

			receiver, err := j.storage.GetUserByID(collaboration.UserID)

			if err != nil {
				return err
			}

			imgURL, err := createImgURL(requester)

			if err != nil {
				return err
			}

			notification := &db.Notification{
				UserID:           collaboration.UserID,
				NotificationType: db.NotificationTypeUserCollaboration,
				Text:             fmt.Sprintf("%s sends you a collaboration message: %s", *requester.FirstName, collaboration.Message),
				EntityType:       "user_collaboration_requests",
				EntityID:         collaboration.ID,
				ChatID:           receiver.ChatID,
				ImageURL:         imgURL,
			}

			created, err := j.storage.CreateNotification(*notification)

			if err != nil {
				return err
			}

			linkToProfile := fmt.Sprintf("https://peatch.pages.dev/users/%d", collaboration.RequesterID)

			if err = j.notifier.SendNotification(created.ChatID, created.Text, imgURL, linkToProfile); err != nil {
				log.Printf("Failed to send notification to user %d", collaboration.UserID)
				return err
			}

			if err = j.storage.UpdateNotificationSentAt(created.ID); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			log.Printf("Notification already sent for user %d and collaboration %d", collaboration.UserID, collaboration.ID)
		}
	}

	return nil
}

func createImgURL(user *db.User) (string, error) {
	// like https://peatch-image-preview.vercel.app/api/image?title=John Doe&subtitle=Product &avatar=https://d262mborv4z66f.cloudfront.net/users/149/KO7uaU43.svg&color=FF8C42&tags=Mentor,17BEBB,e8d3;Founder,FF8C42,eb39;Business Developer,93961F,e992;AI Engineer,685155,f882;Investor,FE5F55,e2eb;Dog Father,685155,f149;Entrepreneur,EF5DA8,e7c8
	base := "https://peatch-image-preview.vercel.app/api/image"

	if user.AvatarURL != nil && user.FirstName != nil && user.LastName != nil && user.Title != nil {
		base = fmt.Sprintf("%s?title=%s %s&subtitle=%s&avatar=https://assets.peatch.io/%s", base, *user.FirstName, *user.LastName, *user.Title, *user.AvatarURL)

		if len(user.Badges) > 0 {
			base = fmt.Sprintf("%s&tags=", base)

			for _, badge := range user.Badges {
				base = fmt.Sprintf("%s%s,%s,%s;", base, badge.Text, badge.Color, badge.Icon)
			}
		}

		return base, nil
	}

	return "", errors.New("user data is not complete")
}
