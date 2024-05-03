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

func (s *service) GetCollaborationByID(userID, id int64) (*db.Collaboration, error) {
	res, err := s.storage.GetCollaborationByID(userID, id)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return nil, terrors.NotFound(err)
	} else if err != nil {
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
	BadgeIDs      []int64 `json:"badge_ids" validate:"dive,min=1,required"`
} // @Name CreateCollaboration

func (cc *CreateCollaboration) toCollaboration() db.Collaboration {
	collab := db.Collaboration{
		OpportunityID: cc.OpportunityID,
		Title:         cc.Title,
		Description:   cc.Description,
		IsPayable:     cc.IsPayable,
		Country:       cc.Country,
		City:          &cc.City,
		CountryCode:   cc.CountryCode,
	}

	if cc.City == "" {
		collab.City = nil
	}

	return collab
}

func (s *service) CreateCollaboration(userID int64, create CreateCollaboration) (*db.Collaboration, error) {
	return s.storage.CreateCollaboration(userID, create.toCollaboration(), create.BadgeIDs)
}

func (s *service) UpdateCollaboration(userID, collabID int64, update CreateCollaboration) error {
	err := s.storage.UpdateCollaboration(userID, collabID, update.toCollaboration(), update.BadgeIDs)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *service) PublishCollaboration(userID int64, collaborationID int64) error {
	err := s.storage.PublishCollaboration(userID, collaborationID)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *service) HideCollaboration(userID int64, collaborationID int64) error {
	err := s.storage.HideCollaboration(userID, collaborationID)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err)
	} else if err != nil {
		return err
	}

	return nil
}

type CreateCollaborationRequest struct {
	Message string `json:"message" validate:"max=1000,required"`
} //@Name CreateCollaborationRequest

func (cr CreateCollaborationRequest) toCollaborationRequest() db.CollaborationRequest {
	return db.CollaborationRequest{
		Message: cr.Message,
	}
}

func (s *service) CreateCollaborationRequest(userID, collaborationID int64, request CreateCollaborationRequest) (*db.CollaborationRequest, error) {
	res, err := s.storage.CreateCollaborationRequest(userID, collaborationID, request.Message)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return nil, terrors.NotFound(err)
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) ShowCollaboration(userID int64, collaborationID int64) error {
	err := s.storage.ShowCollaboration(userID, collaborationID)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *service) FindCollaborationRequest(userID, collabID int64) (*db.CollaborationRequest, error) {
	res, err := s.storage.FindCollaborationRequest(userID, collabID)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, terrors.NotFound(err)
		}

		return nil, err
	}

	return res, nil
}
