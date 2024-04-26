package db

import "time"

type Collaboration struct {
	ID            int64     `json:"id" db:"id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	OpportunityID int64     `json:"opportunity_id" db:"opportunity_id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	IsPayable     bool      `json:"is_payable" db:"is_payable"`
	IsPublished   bool      `json:"is_published" db:"is_published"`
	PublishedAt   string    `json:"published_at" db:"published_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	Country       string    `json:"country" db:"country"`
	City          string    `json:"city" db:"city"`
	CountryCode   string    `json:"country_code" db:"country_code"`
	RequestsCount int       `json:"requests_count" db:"requests_count"`
} // @Name Collaboration

type CollaborationQuery struct {
	Page       int
	Limit      int
	Published  *bool
	OrderBy    CollectionQueryOrder
	SearchTerm string
}

type CollectionQueryOrder string

const (
	CollectionQueryOrderByDate CollectionQueryOrder = "created_at"
)

func (s *storage) ListCollaborations(query CollaborationQuery) ([]Collaboration, error) {
	collaborations := make([]Collaboration, 0)

	queryString := `
		SELECT id, user_id, opportunity_id, title, description, is_payable, is_published, published_at, created_at, updated_at, country, city, country_code, requests_count
		FROM collaborations
		WHERE 1=1
	`

	if query.Published != nil {
		if *query.Published {
			queryString += " AND is_published = true"
		} else {
			queryString += " AND is_published = false"
		}
	}

	if query.SearchTerm != "" {
		queryString += " AND (title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')"
	}

	if query.OrderBy == CollectionQueryOrderByDate {
		queryString += " ORDER BY created_at DESC"
	}

	offset := (query.Page - 1) * query.Limit

	queryString += " OFFSET $2 LIMIT $3"

	err := s.pg.Select(&collaborations, queryString, query.SearchTerm, offset, query.Limit)
	if err != nil {
		return nil, err
	}

	return collaborations, nil
}

func (s *storage) GetCollaborationByID(id int64) (*Collaboration, error) {
	var collaboration Collaboration

	query := `
		SELECT id, user_id, opportunity_id, title, description, is_payable, is_published, published_at, created_at, updated_at, country, city, country_code, requests_count
		FROM collaborations
		WHERE id = $1
	`

	err := s.pg.Get(&collaboration, query, id)
	if err != nil {
		return nil, err
	}

	return &collaboration, nil
}

func (s *storage) CreateCollaboration(collaboration Collaboration) (*Collaboration, error) {
	query := `
		INSERT INTO collaborations (user_id, opportunity_id, title, description, is_payable, country, city, country_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	err := s.pg.QueryRowx(
		query,
		collaboration.UserID,
		collaboration.OpportunityID,
		collaboration.Title,
		collaboration.Description,
		collaboration.IsPayable,
		collaboration.Country,
		collaboration.City,
		collaboration.CountryCode,
	).StructScan(&collaboration)

	if err != nil {
		return nil, err
	}

	return &collaboration, nil
}

func (s *storage) UpdateCollaboration(collaboration Collaboration) (*Collaboration, error) {
	query := `
		UPDATE collaborations
		SET title = $1, description = $2, is_payable = $3, country = $4, city = $5, country_code = $6
		WHERE id = $7
		RETURNING updated_at
	`

	err := s.pg.QueryRowx(
		query,
		collaboration.Title,
		collaboration.Description,
		collaboration.IsPayable,
		collaboration.Country,
		collaboration.City,
		collaboration.CountryCode,
		collaboration.ID,
	).StructScan(&collaboration)

	if err != nil {
		return nil, err
	}

	return &collaboration, nil
}

type CollaborationRequest struct {
	ID              int64     `json:"id" db:"id"`
	CollaborationID int64     `json:"collaboration_id" db:"collaboration_id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	Message         string    `json:"message" db:"message"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	Status          string    `json:"status" db:"status"`
}

func (s *storage) CreateCollaborationRequest(request CollaborationRequest) (*CollaborationRequest, error) {
	query := `
		INSERT INTO collaboration_requests (collaboration_id, user_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := s.pg.QueryRowx(query, request.CollaborationID, request.UserID, request.Message).StructScan(&request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *storage) HideCollaboration(collaborationID int64) error {
	query := `
		UPDATE collaborations
		SET is_published = false, published_at = NULL
		WHERE id = $1
	`

	_, err := s.pg.Exec(query, collaborationID)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) PublishCollaboration(collaborationID int64) error {
	query := `
		UPDATE collaborations
		SET is_published = true, published_at = NOW()
		WHERE id = $1
	`

	_, err := s.pg.Exec(query, collaborationID)
	if err != nil {
		return err
	}

	return nil
}
