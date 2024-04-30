package service

import (
	"errors"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
)

func (s *service) ListCollaborations(query db.CollaborationQuery) ([]db.Collaboration, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	return s.storage.ListCollaborations(query)
}

func (s *service) GetCollaborationByID(id int64) (*db.Collaboration, error) {
	if id == 0 {
		return nil, nil
	}

	res, err := s.storage.GetCollaborationByID(id)

	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			return nil, terrors.NotFound(err)
		}

		return nil, err
	}

	return res, nil
}

type CreateCollaboration struct {
	OpportunityID int64   `json:"opportunity_id" validate:"required"`
	Title         string  `json:"title" validate:"max=255,required"`
	Description   string  `json:"description" validate:"max=1000,required"`
	IsPayable     bool    `json:"is_payable"`
	Country       string  `json:"country" validate:"max=255,required"`
	City          string  `json:"city"`
	CountryCode   string  `json:"country_code" validate:"max=2,required"`
	BadgeIDs      []int64 `json:"badge_ids" validate:"dive,min=1"`
} // @Name CreateCollaboration

func (cc *CreateCollaboration) toCollaboration() db.Collaboration {
	return db.Collaboration{
		OpportunityID: cc.OpportunityID,
		Title:         cc.Title,
		Description:   cc.Description,
		IsPayable:     cc.IsPayable,
		Country:       cc.Country,
		City:          &cc.City,
		CountryCode:   cc.CountryCode,
	}
}

func (s *service) CreateCollaboration(userID int64, create CreateCollaboration) (*db.Collaboration, error) {
	return s.storage.CreateCollaboration(userID, create.toCollaboration(), create.BadgeIDs)
}

func (s *service) UpdateCollaboration(userID int64, update CreateCollaboration) (*db.Collaboration, error) {
	res, err := s.storage.UpdateCollaboration(userID, update.toCollaboration())

	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			return nil, terrors.NotFound(err)
		}

		return nil, err
	}

	return res, nil
}

func (s *service) PublishCollaboration(userID int64, collaborationID int64) error {
	err := s.storage.PublishCollaboration(userID, collaborationID)

	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			return terrors.NotFound(err)
		}

		return err
	}

	return nil
}

func (s *service) HideCollaboration(userID int64, collaborationID int64) error {
	err := s.storage.HideCollaboration(userID, collaborationID)

	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			return terrors.NotFound(err)
		}

		return err
	}

	return nil
}

type CreateCollaborationRequest struct {
	CollaborationID int64  `json:"collaboration_id" validate:"required"`
	Message         string `json:"message" validate:"max=1000"`
}

func (cr CreateCollaborationRequest) toCollaborationRequest() db.CollaborationRequest {
	return db.CollaborationRequest{
		CollaborationID: cr.CollaborationID,
		Message:         cr.Message,
	}
}

func (s *service) CreateCollaborationRequest(userID int64, request CreateCollaborationRequest) (*db.CollaborationRequest, error) {
	res, err := s.storage.CreateCollaborationRequest(userID, request.toCollaborationRequest())

	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			return nil, terrors.NotFound(err)
		}

		return nil, err
	}

	return res, nil
}
