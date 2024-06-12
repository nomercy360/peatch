package service

import (
	"errors"
	"fmt"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
	"regexp"
	"strings"
)

type CreatePostRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	ImageURL    *string `json:"image_url"`
	Country     *string `json:"country"`
	City        *string `json:"city"`
	CountryCode *string `json:"country_code"`
}

func (s *service) CreatePost(userID int64, post CreatePostRequest) (*db.Post, error) {
	imgURL := getImageUrl(post.ImageURL, s.config.CdnURL)

	p := db.Post{
		UserID:      userID,
		Title:       post.Title,
		Description: post.Description,
		ImageURL:    imgURL,
		Country:     post.Country,
		City:        post.City,
		CountryCode: post.CountryCode,
	}

	return s.storage.CreatePost(p)
}

func isValidImagePath(path string) bool {
	// should be [number id of user]/anystring.ext
	re := regexp.MustCompile(`^\d+/.+\.(png|jpg|jpeg|gif)$`)
	return re.MatchString(path)
}

func getImageUrl(url *string, cdnUrl string) *string {
	if url != nil {
		imageURL := *url
		if strings.HasPrefix(imageURL, cdnUrl) {
			return &imageURL
		} else if isValidImagePath(imageURL) {
			imageURL = fmt.Sprintf("%s/%s", cdnUrl, imageURL)
			return &imageURL
		} else {
			return nil
		}
	}

	return nil
}

func (s *service) GetPostByID(uid, id int64) (*db.Post, error) {
	return s.storage.GetPostByID(uid, id)
}

func (s *service) UpdatePost(userID int64, postID int64, updateRequest CreatePostRequest) (*db.Post, error) {
	imgURL := getImageUrl(updateRequest.ImageURL, s.config.CdnURL)

	res, err := s.storage.UpdatePost(userID, postID, db.Post{
		Title:       updateRequest.Title,
		Description: updateRequest.Description,
		ImageURL:    imgURL,
		Country:     updateRequest.Country,
		City:        updateRequest.City,
		CountryCode: updateRequest.CountryCode,
	})

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return nil, terrors.NotFound(err)
	} else if err != nil {
		return nil, terrors.InternalServerError(err)
	}

	return res, nil
}
