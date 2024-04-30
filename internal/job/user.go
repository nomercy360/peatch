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
	ListUserCollaborations(from time.Time) ([]db.UserCollaborationRequest, error)
	UpdateNotificationSentAt(notificationID int64) error
	ListNewCollaborations(from time.Time) ([]db.Collaboration, error)
	FindMatchingUsers(opportunityIDs []int64, badgeIDs []int64) ([]db.User, error)
	ListNewUserProfiles(from time.Time) ([]db.User, error)
}

type notifyJob struct {
	storage       storage
	notifier      notifier
	imgServiceURL string
}

type notifier interface {
	SendNotification(chatID int64, message string, imgUrl string, link string) error
}

func NewNotifyJob(storage storage, notifier notifier, imgServiceURL string) *notifyJob {
	return &notifyJob{
		storage:       storage,
		notifier:      notifier,
		imgServiceURL: imgServiceURL,
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
		userDetails, err := j.storage.GetUserByID(user.ID)

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

		imgURL, err := createImgURL(j.imgServiceURL, userDetails)

		for _, receiver := range receivers {
			_, err := j.storage.SearchNotification(
				receiver.ID,
				db.NotificationTypeUserPublished,
				"users",
				user.ID,
			)

			if err != nil && errors.As(err, &db.ErrNotFound) {

				notification := &db.Notification{
					UserID:           receiver.ID,
					NotificationType: db.NotificationTypeUserPublished,
					Text:             fmt.Sprintf("%s has just published profile", *user.FirstName),
					EntityType:       "users",
					EntityID:         user.ID,
					ChatID:           receiver.ChatID,
					ImageURL:         imgURL,
				}

				created, err := j.storage.CreateNotification(*notification)

				if err != nil {
					return err
				}

				linkToProfile := fmt.Sprintf("https://peatch.pages.dev/users/%d", user.ID)

				if err = j.notifier.SendNotification(created.ChatID, created.Text, imgURL, linkToProfile); err != nil {
					log.Printf("Failed to send notification to user %d", user.ID)
					return err
				}

				if err = j.storage.UpdateNotificationSentAt(created.ID); err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				log.Printf("Notification already sent for user %d and user %d", user.ID, user.ID)
			}
		}
	}

	return nil
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

		imgURL, err := createImgURL(j.imgServiceURL, creator)

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

			imgURL, err := createImgURL(j.imgServiceURL, requester)

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

func createImgURL(baseUrl string, user *db.User) (string, error) {
	// like https://peatch-image-preview.vercel.app/api/image?title=John Doe&subtitle=Product &avatar=https://d262mborv4z66f.cloudfront.net/users/149/KO7uaU43.svg&color=FF8C42&tags=Mentor,17BEBB,e8d3;Founder,FF8C42,eb39;Business Developer,93961F,e992;AI Engineer,685155,f882;Investor,FE5F55,e2eb;Dog Father,685155,f149;Entrepreneur,EF5DA8,e7c8
	baseUrl = fmt.Sprintf("%s/api/image", baseUrl)

	if user.AvatarURL != nil && user.FirstName != nil && user.LastName != nil && user.Title != nil {
		baseUrl = fmt.Sprintf("%s?title=%s %s&subtitle=%s&avatar=https://assets.peatch.io/%s", baseUrl, *user.FirstName, *user.LastName, *user.Title, *user.AvatarURL)

		if len(user.Badges) > 0 {
			baseUrl = fmt.Sprintf("%s&tags=", baseUrl)

			for _, badge := range user.Badges {
				baseUrl = fmt.Sprintf("%s%s,%s,%s;", baseUrl, badge.Text, badge.Color, badge.Icon)
			}
		}

		return baseUrl, nil
	}

	return "", errors.New("user data is not complete")
}
