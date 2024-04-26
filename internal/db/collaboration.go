package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type Collaboration struct {
	ID            int64       `json:"id" db:"id"`
	UserID        int64       `json:"user_id" db:"user_id"`
	OpportunityID int64       `json:"opportunity_id" db:"opportunity_id"`
	Title         string      `json:"title" db:"title"`
	Description   string      `json:"description" db:"description"`
	IsPayable     bool        `json:"is_payable" db:"is_payable"`
	PublishedAt   *string     `json:"published_at" db:"published_at"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	Country       string      `json:"country" db:"country"`
	City          string      `json:"city" db:"city"`
	CountryCode   string      `json:"country_code" db:"country_code"`
	RequestsCount int         `json:"requests_count" db:"requests_count"`
	Badges        []Badge     `json:"badges" db:"-"`
	Opportunity   Opportunity `json:"opportunity" db:"-"`
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

func scanCollaborationWithOpportunity(rows *sqlx.Rows) (Collaboration, error) {
	var collaboration Collaboration
	var opportunity Opportunity

	if err := rows.Scan(
		&collaboration.ID, &collaboration.UserID, &collaboration.OpportunityID,
		&collaboration.Title, &collaboration.Description, &collaboration.IsPayable,
		&collaboration.PublishedAt, &collaboration.CreatedAt, &collaboration.UpdatedAt,
		&collaboration.Country, &collaboration.City, &collaboration.CountryCode,
		&collaboration.RequestsCount, &opportunity.ID, &opportunity.Text,
		&opportunity.Description, &opportunity.Icon, &opportunity.Color, &opportunity.CreatedAt,
	); err != nil {
		return Collaboration{}, err
	}

	collaboration.Opportunity = opportunity

	return collaboration, nil
}

func (s *storage) ListCollaborations(query CollaborationQuery) ([]Collaboration, error) {
	collaborations := make([]Collaboration, 0)
	var args []interface{}
	paramIndex := 1

	queryString := `
		SELECT c.id, c.user_id, c.opportunity_id, c.title, c.description, c.is_payable, c.published_at, c.created_at, c.updated_at, c.country, c.city, c.country_code, c.requests_count,
			o.id, o.text, o.description, o.icon, o.color, o.created_at
		FROM collaborations c
		JOIN opportunities o ON c.opportunity_id = o.id
		WHERE published_at IS NOT NULL
	`

	if query.SearchTerm != "" {
		args = append(args, query.SearchTerm)
		queryString += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", paramIndex, paramIndex)
		paramIndex++
	}

	if query.OrderBy == CollectionQueryOrderByDate {
		queryString += " ORDER BY created_at DESC"
	}

	offset := (query.Page - 1) * query.Limit

	queryString += fmt.Sprintf(" OFFSET $%d LIMIT $%d", paramIndex, paramIndex+1)
	args = append(args, offset, query.Limit)

	rows, err := s.pg.Queryx(queryString, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		collaboration, err := scanCollaborationWithOpportunity(rows)
		if err != nil {
			return nil, err
		}

		collaborations = append(collaborations, collaboration)
	}

	return collaborations, nil
}

func (s *storage) GetCollaborationByID(id int64) (*Collaboration, error) {
	var collaboration Collaboration

	query := `
        SELECT 
            c.id, c.user_id, c.opportunity_id, c.title, c.description, c.is_payable, c.published_at, c.created_at, c.updated_at, c.country, c.city, c.country_code, c.requests_count,
            o.id, o.text, o.description, o.icon, o.color, o.created_at
        FROM collaborations c
        LEFT JOIN opportunities o ON c.opportunity_id = o.id
        WHERE c.id = $1 AND c.published_at IS NOT NULL
    `

	rows, err := s.pg.Queryx(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		collaboration, err = scanCollaborationWithOpportunity(rows)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrNotFound
	}

	// fetch badges
	query = `
		SELECT b.id, b.text, b.icon, b.color, b.created_at
		FROM badges b
		JOIN collaboration_badges cb ON b.id = cb.badge_id
		WHERE cb.collaboration_id = $1
	`

	err = s.pg.Select(&collaboration.Badges, query, id)

	if err != nil {
		return nil, err
	}

	return &collaboration, nil
}

func (s *storage) CreateCollaboration(userID int64, collaboration Collaboration, badges []int64) (*Collaboration, error) {
	var res Collaboration

	query := `
		INSERT INTO collaborations (user_id, opportunity_id, title, description, is_payable, country, city, country_code, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, user_id, opportunity_id, title, description, is_payable, published_at, created_at, updated_at, country, city, country_code, requests_count
	`

	err := s.pg.QueryRowx(
		query,
		userID,
		collaboration.OpportunityID,
		collaboration.Title,
		collaboration.Description,
		collaboration.IsPayable,
		collaboration.Country,
		collaboration.City,
		collaboration.CountryCode,
	).StructScan(&res)

	if err != nil {
		return nil, err
	}

	if len(badges) > 0 {
		var valueStrings []string
		var valueArgs []interface{}
		for _, badge := range badges {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, res.ID, badge)
		}

		stmt := `INSERT INTO collaboration_badges (collaboration_id, badge_id) VALUES ` + strings.Join(valueStrings, ", ")
		stmt = s.pg.Rebind(stmt)

		_, err = s.pg.Exec(stmt, valueArgs...)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *storage) UpdateCollaboration(userID int64, collaboration Collaboration) (*Collaboration, error) {
	query := `
		UPDATE collaborations
		SET title = $1, description = $2, is_payable = $3, country = $4, city = $5, country_code = $6
		WHERE id = $7 AND user_id = $8
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
		userID,
	).StructScan(&collaboration)

	if err != nil {
		if IsNoRowsError(err) {
			return nil, ErrNotFound
		}

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

func (s *storage) CreateCollaborationRequest(userID int64, request CollaborationRequest) (*CollaborationRequest, error) {
	query := `
		INSERT INTO collaboration_requests (collaboration_id, user_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := s.pg.QueryRowx(query, request.CollaborationID, userID, request.Message).StructScan(&request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *storage) HideCollaboration(userID int64, collaborationID int64) error {
	query := `
		UPDATE collaborations
		SET published_at = NULL
		WHERE id = $1 AND user_id = $2
	`

	res, err := s.pg.Exec(query, collaborationID, userID)
	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *storage) PublishCollaboration(userID int64, collaborationID int64) error {
	query := `
		UPDATE collaborations
		SET published_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	res, err := s.pg.Exec(query, collaborationID, userID)
	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
