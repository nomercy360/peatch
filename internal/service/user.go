package service

import (
	"github.com/peatch-io/peatch/internal/db"
)

func (s *service) ListUsers(query db.UserQuery) ([]db.User, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	var published bool
	query.Published = &published

	return s.storage.ListUsers(query)
}

func (s *service) GetUserByChatID(chatID int64) (*db.User, error) {
	if chatID == 0 {
		return nil, nil
	}

	return s.storage.GetUserByChatID(chatID)
}

func (s *service) CreateUser(user db.User) (*db.User, error) {
	return s.storage.CreateUser(user)
}

func (s *service) GetUserByID(id int64) (*db.User, error) {
	if id == 0 {
		return nil, nil
	}

	return s.storage.GetUserByID(id)
}

func (s *service) UpdateUser(user db.User) (*db.User, error) {
	return s.storage.UpdateUser(user)
}

func (s *service) FollowUser(userID, followerID int64) error {
	return s.storage.FollowUser(userID, followerID)
}

func (s *service) UnfollowUser(userID, followerID int64) error {
	return s.storage.UnfollowUser(userID, followerID)
}

func (s *service) PublishUser(userID int64) error {
	return s.storage.PublishUser(userID)
}

func (s *service) HideUser(userID int64) error {
	return s.storage.HideUser(userID)
}
