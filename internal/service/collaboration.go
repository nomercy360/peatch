package service

import "github.com/peatch-io/peatch/internal/db"

func (s *service) ListCollaborations(query db.CollaborationQuery) ([]db.Collaboration, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	var published bool
	query.Published = &published

	return s.storage.ListCollaborations(query)
}

func (s *service) GetCollaborationByID(id int64) (*db.Collaboration, error) {
	if id == 0 {
		return nil, nil
	}

	return s.storage.GetCollaborationByID(id)
}

func (s *service) CreateCollaboration(collaboration db.Collaboration) (*db.Collaboration, error) {
	return s.storage.CreateCollaboration(collaboration)
}

func (s *service) UpdateCollaboration(collaboration db.Collaboration) (*db.Collaboration, error) {
	return s.storage.UpdateCollaboration(collaboration)
}

func (s *service) PublishCollaboration(collaborationID int64) error {
	return s.storage.PublishCollaboration(collaborationID)
}

func (s *service) HideCollaboration(collaborationID int64) error {
	return s.storage.HideCollaboration(collaborationID)
}

func (s *service) CreateCollaborationRequest(request db.CollaborationRequest) (*db.CollaborationRequest, error) {
	return s.storage.CreateCollaborationRequest(request)
}
