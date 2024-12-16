package service

import (
	"github.com/peatch-io/peatch/internal/db"
)

type FeedQuery struct {
	Page   int
	Limit  int
	Search string
}

type FeedItem struct {
	Type string  `json:"type"`
	Data db.Base `json:"data"`
}

func (s *service) GetFeed(uid int64, query FeedQuery) ([]FeedItem, error) {
	items := make([]FeedItem, 0)

	if query.Limit <= 0 {
		query.Limit = 40
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	res, err := s.storage.ListUsers(db.UserQuery{
		UserID: uid,
		Page:   query.Page,
		Limit:  query.Limit,
		Search: query.Search,
	})

	if err != nil {
		return nil, err
	}

	for _, user := range res {
		feedItem := FeedItem{
			Type: "user",
			Data: user.ToUserProfile(),
		}

		items = append(items, feedItem)
	}

	// fetch collaborations
	c, err := s.storage.ListCollaborations(db.CollaborationQuery{
		Page:   query.Page,
		Limit:  query.Limit,
		Search: query.Search,
		UserID: uid,
	})

	if err != nil {
		return nil, err
	}

	for _, collab := range c {
		feedItem := FeedItem{
			Type: "collaboration",
			Data: collab,
		}

		items = append(items, feedItem)
	}

	posts := make([]db.Post, 0)

	// fetch posts
	posts, err = s.storage.GetPosts(db.PostQuery{
		Search: query.Search,
		Page:   query.Page,
		Limit:  query.Limit,
		UserID: uid,
	})

	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		feedItem := FeedItem{
			Type: "post",
			Data: post,
		}

		items = append(items, feedItem)
	}

	return items, nil
}
