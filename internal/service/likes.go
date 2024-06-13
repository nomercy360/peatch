package service

import (
	"errors"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
)

type LikeRequest struct {
	ContentID   int64  `json:"content_id" validate:"required"`
	ContentType string `json:"content_type" validate:"required,oneof=post user collaboration"`
}

func (s *service) IncreaseLikeCount(userID int64, req LikeRequest) error {
	err := s.storage.IncreaseLikeCount(userID, req.ContentID, req.ContentType)

	if err != nil && errors.Is(err, db.ErrAlreadyExists) {
		return terrors.BadRequest(err)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *service) DecreaseLikeCount(userID int64, req LikeRequest) error {
	err := s.storage.DecreaseLikeCount(userID, req.ContentID, req.ContentType)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err)
	} else if err != nil {
		return err
	}

	return nil
}
