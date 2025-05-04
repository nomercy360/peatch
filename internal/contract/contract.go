package contract

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peatch-io/peatch/internal/db"
	"time"
)

type OkResponse struct {
	Message string `json:"message"`
} // @Name OkResponse

type ErrorResponse struct {
	Error string `json:"error"`
} // @Name ErrorResponse

type StatusResponse struct {
	Success bool `json:"success"`
} // @Name StatusResponse

type AuthTelegramRequest struct {
	Query string `json:"query"`
} // @Name AuthTelegramRequest

func (a AuthTelegramRequest) Validate() error {
	if a.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	return nil
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UID    string `json:"uid"`
	ChatID int64  `json:"chat_id"`
	Lang   string `json:"lang"`
}

type UpdateUserRequest struct {
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	LocationID     string   `json:"location_id"`
	BadgeIDs       []string `json:"badge_ids"`
	OpportunityIDs []string `json:"opportunity_ids"`
} // @Name UpdateUserRequest

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
} // @Name Location

type CityResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	CountryCode string   `json:"country_code"`
	CountryName string   `json:"country_name"`
	Location    Location `json:"location"`
} // @Name CityResponse

func ToCityResponse(city db.City) CityResponse {
	loc := Location{}
	if len(city.Geo.Coordinates) == 2 {
		loc.Latitude = city.Geo.Coordinates[1]
		loc.Longitude = city.Geo.Coordinates[0]
	}
	return CityResponse{
		ID:          city.ID,
		Name:        city.Name,
		CountryCode: city.CountryCode,
		CountryName: city.CountryName,
		Location:    loc,
	}
}

func (r UpdateUserRequest) Validate() error {
	if r.FirstName == "" {
		return fmt.Errorf("first_name is required")
	}
	if r.LastName == "" {
		return fmt.Errorf("last_name is required")
	}
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(r.Title) > 255 {
		return fmt.Errorf("title must not exceed 255 characters")
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	if len(r.Description) > 1000 {
		return fmt.Errorf("description must not exceed 1000 characters")
	}
	if r.LocationID == "" {
		return fmt.Errorf("location_id is required")
	}
	if len(r.BadgeIDs) == 0 {
		return fmt.Errorf("badge_ids is required")
	}
	if len(r.OpportunityIDs) == 0 {
		return fmt.Errorf("opportunity_ids is required")
	}
	return nil
}

type CreateCollaboration struct {
	OpportunityID string   `json:"opportunity_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	IsPayable     bool     `json:"is_payable"`
	LocationID    string   `json:"location_id"`
	BadgeIDs      []string `json:"badge_ids"`
} // @Name CreateCollaboration

func (r CreateCollaboration) Validate() error {
	if r.OpportunityID == "" {
		return fmt.Errorf("opportunity_id is required")
	}
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(r.Title) > 255 {
		return fmt.Errorf("title must not exceed 255 characters")
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	if len(r.Description) > 1000 {
		return fmt.Errorf("description must not exceed 1000 characters")
	}
	if r.LocationID == "" {
		return fmt.Errorf("location_id is required")
	}
	if len(r.BadgeIDs) == 0 {
		return fmt.Errorf("badge_ids must contain at least one element")
	}
	return nil
}

type CreateBadgeRequest struct {
	Text  string `json:"text"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
} // @Name CreateBadgeRequest

func (r CreateBadgeRequest) Validate() error {
	if r.Text == "" {
		return fmt.Errorf("text is required")
	}
	if r.Icon == "" {
		return fmt.Errorf("icon is required")
	}
	if len(r.Icon) != 4 {
		return fmt.Errorf("icon must be exactly 4 characters and hexadecimal")
	}
	if r.Color == "" {
		return fmt.Errorf("color is required")
	}
	if len(r.Color) != 6 {
		return fmt.Errorf("color must be exactly 6 characters and hexadecimal")
	}
	return nil
}

type OpportunityResponse struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
} // @Name OpportunityResponse

type UserResponse struct {
	ID                 string                `json:"id"`
	ChatID             int64                 `json:"chat_id"`
	Username           string                `json:"username"`
	FirstName          *string               `json:"first_name"`
	LastName           *string               `json:"last_name"`
	LanguageCode       db.LanguageCode       `json:"language_code"`
	AvatarURL          *string               `json:"avatar_url"`
	Title              *string               `json:"title"`
	Description        *string               `json:"description"`
	Location           *CityResponse         `json:"location"`
	Badges             []BadgeResponse       `json:"badges"`
	Opportunities      []OpportunityResponse `json:"opportunities"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
	LastActiveAt       time.Time             `json:"last_active_at"`
	VerificationStatus db.VerificationStatus `json:"verification_status"`
} // @Name UserResponse

func ToUserResponse(user db.User) UserResponse {
	location := &CityResponse{}
	if user.Location != nil {
		location = &CityResponse{
			ID:          user.Location.ID,
			Name:        user.Location.Name,
			CountryCode: user.Location.CountryCode,
			CountryName: user.Location.CountryName,
			Location: Location{
				Latitude:  user.Location.Geo.Coordinates[1],
				Longitude: user.Location.Geo.Coordinates[0],
			},
		}
	}
	return UserResponse{
		ID:                 user.ID,
		ChatID:             user.ChatID,
		Username:           user.Username,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		LanguageCode:       user.LanguageCode,
		AvatarURL:          user.AvatarURL,
		Title:              user.Title,
		Description:        user.Description,
		Location:           location,
		LastActiveAt:       user.LastActiveAt,
		Badges:             ToBadgeResponseList(user.Badges),
		Opportunities:      ToOpportunityResponseList(user.Opportunities, user.LanguageCode),
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
		VerificationStatus: user.VerificationStatus,
	}
}

type UserProfileResponse struct {
	ID            string                `json:"id"`
	FirstName     string                `json:"first_name"`
	LastName      string                `json:"last_name"`
	AvatarURL     string                `json:"avatar_url"`
	Title         string                `json:"title"`
	Description   string                `json:"description"`
	Location      CityResponse          `json:"location"`
	IsFollowing   bool                  `json:"is_following"`
	Badges        []BadgeResponse       `json:"badges"`
	Opportunities []OpportunityResponse `json:"opportunities"`
	LastActiveAt  time.Time             `json:"last_active_at"`
} // @Name UserProfileResponse

func ToUserProfile(user db.User) UserProfileResponse {
	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}
	avatarURL := ""
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}
	title := ""
	if user.Title != nil {
		title = *user.Title
	}
	description := ""
	if user.Description != nil {
		description = *user.Description
	}
	location := CityResponse{}
	if user.Location != nil {
		location = ToCityResponse(*user.Location)
	}
	return UserProfileResponse{
		ID:            user.ID,
		FirstName:     firstName,
		LastName:      lastName,
		AvatarURL:     avatarURL,
		Title:         title,
		Description:   description,
		Location:      location,
		IsFollowing:   user.IsFollowing,
		Badges:        ToBadgeResponseList(user.Badges),
		Opportunities: ToOpportunityResponseList(user.Opportunities, user.LanguageCode),
		LastActiveAt:  user.LastActiveAt,
	}
}

type BadgeResponse struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
} // @Name BadgeResponse

func ToBadgeResponse(badge db.Badge) BadgeResponse {
	return BadgeResponse{
		ID:    badge.ID,
		Text:  badge.Text,
		Icon:  badge.Icon,
		Color: badge.Color,
	}
}

func ToBadgeResponseList(badges []db.Badge) []BadgeResponse {
	badgeResponses := make([]BadgeResponse, len(badges))
	for i, badge := range badges {
		badgeResponses[i] = ToBadgeResponse(badge)
	}
	return badgeResponses
}

func ToOpportunityResponseList(opportunities []db.Opportunity, lang db.LanguageCode) []OpportunityResponse {
	opportunityResponses := make([]OpportunityResponse, len(opportunities))
	for i, opportunity := range opportunities {
		opportunityResponses[i] = OpportunityResponse{
			ID:    opportunity.ID,
			Icon:  opportunity.Icon,
			Color: opportunity.Color,
		}

		if lang == db.LanguageRU {
			opportunityResponses[i].Text = opportunity.TextRU
			opportunityResponses[i].Description = opportunity.DescriptionRU
		} else {
			opportunityResponses[i].Text = opportunity.Text
			opportunityResponses[i].Description = opportunity.Description
		}
	}
	return opportunityResponses
}

type CollaborationResponse struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	IsPayable   bool                `json:"is_payable"`
	Badges      []BadgeResponse     `json:"badges"`
	Opportunity OpportunityResponse `json:"opportunity"`
	Location    CityResponse        `json:"location"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	User        UserProfileResponse `json:"user"`
} // @Name CollaborationResponse

func ToCollaborationResponse(collab db.Collaboration) CollaborationResponse {
	return CollaborationResponse{
		ID:          collab.ID,
		UserID:      collab.UserID,
		Title:       collab.Title,
		Description: collab.Description,
		IsPayable:   collab.IsPayable,
		Badges:      ToBadgeResponseList(collab.Badges),
		Opportunity: ToOpportunityResponseList([]db.Opportunity{collab.Opportunity}, db.LanguageEN)[0],
		Location:    ToCityResponse(collab.Location),
		CreatedAt:   collab.CreatedAt,
		UpdatedAt:   collab.UpdatedAt,
		User:        ToUserProfile(collab.User),
	}
}
