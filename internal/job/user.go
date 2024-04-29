package job

import "github.com/peatch-io/peatch/internal/db"

type storage interface {
	GetUserByID(id int64) (*db.User, error)
	SaveNotification(notification db.Notification) (*db.Notification, error)
	GetLastSentNotification(userID int64) (*db.Notification, error)
}

type notifyJob struct {
	storage  storage
	notifier notifier
}

type notifier interface {
	SendNotification(userID int64, message string) error
}

func NewNotifyJob(storage storage, notifier notifier) *notifyJob {
	return &notifyJob{
		storage:  storage,
		notifier: notifier,
	}
}

func (j *notifyJob) UserRegistrationJob() error {
	user, err := j.storage.GetUserByID(1)

	if err != nil {
		return err
	}

	return j.notifier.SendNotification(user.ID, "Welcome to our platform!")
}
