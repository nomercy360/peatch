package interfaces

import (
	"github.com/peatch-io/peatch/internal/db"
)

// NotificationService defines the interface for notification services
type NotificationService interface {
	NotifyUserVerified(user db.User) error
	NotifyCollaborationVerified(collab db.Collaboration) error
	NotifyUserVerificationDenied(user db.User) error
	NotifyCollaborationVerificationDenied(collab db.Collaboration) error
	NotifyNewPendingUser(user db.User) error
	NotifyNewPendingCollaboration(user db.User, collab db.Collaboration) error
	NotifyUserFollow(follower db.User, followee db.User) error
	NotifyCollabInterest(user db.User, collab db.Collaboration) error
	SendCollaborationToCommunityChatWithImage(collab db.Collaboration) error
}
