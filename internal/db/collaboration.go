package db

import (
	"fmt"
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
	PublishedAt   *time.Time  `json:"published_at" db:"published_at"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"-" db:"updated_at"`
	Country       string      `json:"country" db:"country"`
	City          *string     `json:"city" db:"city"`
	CountryCode   string      `json:"country_code" db:"country_code"`
	HiddenAt      *time.Time  `json:"hidden_at" db:"hidden_at"`
	RequestsCount int         `json:"-" db:"requests_count"`
	Opportunity   Opportunity `json:"opportunity" db:"opportunity"`
	User          UserProfile `json:"user" db:"user"`
	Badges        BadgeSlice  `json:"badges,omitempty" db:"badges"`
	LikesCount    int         `json:"likes_count" db:"likes_count"`
	IsLiked       bool        `json:"is_liked" db:"is_liked"`
} // @Name Collaboration

func (c *Collaboration) GetLocation() string {
	if c.City != nil {
		return fmt.Sprintf("%s, %s", *c.City, c.Country)
	}

	return c.Country
}

type CollaborationQuery struct {
	Page    int
	Limit   int
	Search  string
	From    *time.Time
	UserID  int64
	Visible bool
}

func (s *storage) ListCollaborations(params CollaborationQuery) ([]Collaboration, error) {
	res := make([]Collaboration, 0)
	query := `
        SELECT c.*,
			to_jsonb(o) as opportunity,
			to_jsonb(u) as "user",
			exists(SELECT 1 FROM likes l WHERE l.content_id = c.id AND l.user_id = $1 AND l.content_type = 'collaboration') as is_liked
        FROM collaborations c
        LEFT JOIN opportunities o ON c.opportunity_id = o.id
		LEFT JOIN users u ON c.user_id = u.id
        WHERE c.published_at IS NOT NULL AND c.hidden_at IS NULL
    `

	paramIndex := 2
	args := []interface{}{params.UserID}

	var whereClauses []string

	if params.Search != "" {
		searchClause := fmt.Sprintf("AND (c.title ILIKE $%d OR c.description ILIKE $%d)", paramIndex, paramIndex)
		args = append(args, "%"+params.Search+"%")
		whereClauses = append(whereClauses, searchClause)
		paramIndex++
	}

	if params.From != nil {
		fromClause := fmt.Sprintf("AND c.created_at >= $%d", paramIndex)
		args = append(args, *params.From)
		whereClauses = append(whereClauses, fromClause)
		paramIndex++
	}

	query = query + strings.Join(whereClauses, " ")
	query += fmt.Sprintf(" ORDER BY c.created_at DESC")
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	err := s.pg.Select(&res, query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *storage) GetCollaborationByID(userID, id int64) (*Collaboration, error) {
	var collaboration Collaboration
	var args []interface{}

	query := `
        SELECT 
            c.*,
			to_jsonb(o) as opportunity,
			to_jsonb(u) as "user",
			json_agg(distinct to_jsonb(b)) as badges
		FROM collaborations c
		LEFT JOIN opportunities o ON c.opportunity_id = o.id
		LEFT JOIN users u ON c.user_id = u.id
		LEFT JOIN collaboration_badges cb ON c.id = cb.collaboration_id
		LEFT JOIN badges b ON cb.badge_id = b.id
		WHERE c.id = $1
	`

	args = append(args, id)

	if userID != 0 {
		query += " AND (c.user_id = $2 OR (c.published_at IS NOT NULL AND c.hidden_at IS NULL))"
		args = append(args, userID)
	}

	query += fmt.Sprintf(" GROUP BY c.id, o.id, u.id")

	err := s.pg.Get(&collaboration, query, args...)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &collaboration, nil
}

func (s *storage) CreateCollaboration(userID int64, collaboration Collaboration, badges []int64) (*Collaboration, error) {
	var res Collaboration

	tx, err := s.pg.Beginx()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO collaborations (user_id, opportunity_id, title, description, is_payable, country, city, country_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, opportunity_id, title, description, is_payable, published_at, hidden_at, created_at, updated_at, country, city, country_code, requests_count
	`

	if err := tx.QueryRowx(
		query,
		userID,
		collaboration.OpportunityID,
		collaboration.Title,
		collaboration.Description,
		collaboration.IsPayable,
		collaboration.Country,
		collaboration.City,
		collaboration.CountryCode,
	).StructScan(&res); err != nil {
		tx.Rollback()
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
		stmt = tx.Rebind(stmt)

		if _, err := tx.Exec(stmt, valueArgs...); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *storage) UpdateCollaboration(userID, collabID int64, collaboration Collaboration, badges []int64) error {
	tx, err := s.pg.Beginx()
	if err != nil {
		return err
	}

	query := `
		UPDATE collaborations
		SET title = $1, description = $2, is_payable = $3, country = $4, city = $5, country_code = $6, updated_at = NOW(), opportunity_id = $7
		WHERE id = $8 AND user_id = $9
		RETURNING updated_at
	`

	err = tx.QueryRowx(
		query,
		collaboration.Title,
		collaboration.Description,
		collaboration.IsPayable,
		collaboration.Country,
		collaboration.City,
		collaboration.CountryCode,
		collaboration.OpportunityID,
		collabID,
		userID,
	).StructScan(&collaboration)

	if err != nil && IsNoRowsError(err) {
		return ErrNotFound
	} else if err != nil {
		tx.Rollback()
		return err
	}

	if len(badges) > 0 {
		_, err = tx.Exec("DELETE FROM collaboration_badges WHERE collaboration_id = $1", collabID)
		if err != nil {
			tx.Rollback()
			return err
		}

		var valueStrings []string
		var valueArgs []interface{}
		for _, badge := range badges {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, collabID, badge)
		}

		stmt := `INSERT INTO collaboration_badges (collaboration_id, badge_id) VALUES ` + strings.Join(valueStrings, ", ")
		stmt = tx.Rebind(stmt)

		if _, err := tx.Exec(stmt, valueArgs...); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

type CollaborationRequest struct {
	ID              int64     `json:"id" db:"id"`
	CollaborationID int64     `json:"collaboration_id" db:"collaboration_id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	Message         string    `json:"message" db:"message"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	Status          string    `json:"status" db:"status"`
} // @Name CollaborationRequest

func (s *storage) CreateCollaborationRequest(userID, collaborationID int64, message string) (*CollaborationRequest, error) {
	var request CollaborationRequest

	query := `
		INSERT INTO collaboration_requests (collaboration_id, user_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, collaboration_id, user_id, message, created_at, updated_at, status
	`

	err := s.pg.QueryRowx(query, collaborationID, userID, message).StructScan(&request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *storage) HideCollaboration(userID int64, collaborationID int64) error {
	query := `
		UPDATE collaborations
		SET hidden_at = NOW()
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

func (s *storage) ShowCollaboration(userID int64, collaborationID int64) error {
	query := `
		UPDATE collaborations
		SET hidden_at = NULL
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

func (s *storage) ListCollaborationRequests(from time.Time) ([]CollaborationRequest, error) {
	requests := make([]CollaborationRequest, 0)

	query := `
		SELECT id, collaboration_id, user_id, message, created_at, updated_at, status
		FROM collaboration_requests
		WHERE created_at >= $1
	`

	err := s.pg.Select(&requests, query, from)

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (s *storage) FindCollaborationRequest(userID, collabID int64) (*CollaborationRequest, error) {
	var request CollaborationRequest

	query := `
		SELECT id, collaboration_id, user_id, message, created_at, updated_at, status
		FROM collaboration_requests
		WHERE collaboration_id = $1 AND user_id = $2
	`

	err := s.pg.Get(&request, query, collabID, userID)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
