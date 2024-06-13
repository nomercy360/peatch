package db

import (
	"fmt"
	"strings"
	"time"
)

type Post struct {
	ID          int64       `json:"id" db:"id"`
	UserID      int64       `json:"user_id" db:"user_id"`
	Title       string      `json:"title" db:"title"`
	Description string      `json:"description" db:"description"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
	HiddenAt    *time.Time  `json:"hidden_at" db:"hidden_at"`
	ImageURL    *string     `json:"image_url" db:"image_url"`
	Country     *string     `json:"country" db:"country"`
	City        *string     `json:"city" db:"city"`
	CountryCode *string     `json:"country_code" db:"country_code"`
	User        UserProfile `json:"user" db:"user"`
	LikesCount  int         `json:"likes_count" db:"likes_count"`
	IsLiked     bool        `json:"is_liked" db:"is_liked"`
} //@Name Post

func (p Post) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (s *storage) CreatePost(post Post) (*Post, error) {
	query := `
		INSERT INTO posts (user_id, title, description, image_url, country, city, country_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, title, description, created_at, updated_at, image_url, country, city, country_code
	`

	err := s.pg.Get(&post, query, post.UserID, post.Title, post.Description, post.ImageURL, post.Country, post.City, post.CountryCode)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *storage) GetPostByID(uid, id int64) (*Post, error) {
	var post Post

	query := `
		SELECT p.id, p.user_id, p.title, p.description, p.created_at, p.updated_at, p.image_url, p.country, p.city, p.country_code,
		       p.likes_count,
		       to_json(u) as "user"
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.id = $1 AND (p.hidden_at IS NULL OR u.id = $2)
	`

	err := s.pg.Get(&post, query, id, uid)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

type PostQuery struct {
	Page   int
	Limit  int
	Search string
	UserID int64
}

func (s *storage) GetPosts(query PostQuery) ([]Post, error) {
	var posts []Post

	q := `
		SELECT p.id, p.user_id, p.title, p.description, p.created_at, p.updated_at, p.image_url, p.country, p.city, p.country_code,
		       p.likes_count,
		       exists(SELECT 1 FROM likes l WHERE l.content_id = p.id AND l.content_type = 'post' AND l.user_id = $1) as is_liked,
		       to_json(u) as "user"
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE 1=1
	`

	paramIndex := 2
	args := []interface{}{query.UserID}
	var whereClauses []string

	if query.Search != "" {
		searchClause := fmt.Sprintf("AND (p.title ILIKE $%d OR p.description ILIKE $%d)", paramIndex, paramIndex)
		args = append(args, "%"+query.Search+"%")
		whereClauses = append(whereClauses, searchClause)
		paramIndex++
	}

	q = q + strings.Join(whereClauses, " ")
	q += fmt.Sprintf(" ORDER BY p.created_at DESC")
	q += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)

	offset := (query.Page - 1) * query.Limit
	args = append(args, query.Limit, offset)

	err := s.pg.Select(&posts, q, args...)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *storage) UpdatePost(uid, postID int64, post Post) (*Post, error) {
	query := `
		UPDATE posts
		SET title = $1, description = $2, image_url = $3, country = $4, city = $5, country_code = $6
		WHERE id = $7 AND user_id = $8
		RETURNING id, user_id, title, description, created_at, updated_at, image_url, country, city, country_code
	`

	err := s.pg.Get(&post, query, post.Title, post.Description, post.ImageURL, post.Country, post.City, post.CountryCode, postID, uid)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
