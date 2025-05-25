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

type BotBlockedResponse struct {
	Status   string `json:"status"`
	Username string `json:"username"`
	Message  string `json:"message"`
} // @Name BotBlockedResponse

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

type Link struct {
	URL   string `json:"url"`
	Label string `json:"label"`
	Type  string `json:"type"`
	Order int    `json:"order"`
} // @Name Link

type UpdateUserRequest struct {
	Name           string   `json:"name"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	LocationID     string   `json:"location_id"`
	BadgeIDs       []string `json:"badge_ids"`
	OpportunityIDs []string `json:"opportunity_ids"`
	Links          []Link   `json:"links"`
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
	if r.Name == "" {
		return fmt.Errorf("name is required")
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

	// Validate links
	for i, link := range r.Links {
		if link.URL == "" {
			return fmt.Errorf("link %d: url is required", i+1)
		}
		if link.Label == "" {
			return fmt.Errorf("link %d: label is required", i+1)
		}
		if len(link.Label) > 100 {
			return fmt.Errorf("link %d: label must not exceed 100 characters", i+1)
		}
		if link.Type == "" {
			return fmt.Errorf("link %d: type is required", i+1)
		}
	}

	return nil
}

type UpdateUserLinksRequest struct {
	Links []Link `json:"links"`
} // @Name UpdateUserLinksRequest

func (r UpdateUserLinksRequest) Validate() error {
	// Validate links
	for i, link := range r.Links {
		if link.URL == "" {
			return fmt.Errorf("link %d: url is required", i+1)
		}
		if link.Label == "" {
			return fmt.Errorf("link %d: label is required", i+1)
		}
		if len(link.Label) > 100 {
			return fmt.Errorf("link %d: label must not exceed 100 characters", i+1)
		}
		if link.Type == "" {
			return fmt.Errorf("link %d: type is required", i+1)
		}
	}
	return nil
}

type CreateCollaboration struct {
	OpportunityID string   `json:"opportunity_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	IsPayable     bool     `json:"is_payable"`
	LocationID    *string  `json:"location_id"`
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
	if r.LocationID != nil && *r.LocationID == "" {
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
	Name               *string               `json:"name"`
	LanguageCode       db.LanguageCode       `json:"language_code"`
	AvatarURL          *string               `json:"avatar_url"`
	Title              *string               `json:"title"`
	Description        *string               `json:"description"`
	Location           *CityResponse         `json:"location"`
	Badges             []BadgeResponse       `json:"badges"`
	Opportunities      []OpportunityResponse `json:"opportunities"`
	Links              []Link                `json:"links"`
	HiddenAt           *time.Time            `json:"hidden_at"`
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
		Name:               user.Name,
		LanguageCode:       user.LanguageCode,
		AvatarURL:          user.AvatarURL,
		Title:              user.Title,
		Description:        user.Description,
		Location:           location,
		LastActiveAt:       user.LastActiveAt,
		Badges:             ToBadgeResponseList(user.Badges),
		Opportunities:      ToOpportunityResponseList(user.Opportunities, user.LanguageCode),
		Links:              ToLinkResponseList(user.Links),
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
		VerificationStatus: user.VerificationStatus,
		HiddenAt:           user.HiddenAt,
	}
}

type UserProfileResponse struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	AvatarURL     string                `json:"avatar_url"`
	Title         string                `json:"title,omitempty"`
	Description   string                `json:"description,omitempty"`
	Location      CityResponse          `json:"location,omitempty"`
	IsFollowing   bool                  `json:"is_following"`
	Badges        []BadgeResponse       `json:"badges"`
	Opportunities []OpportunityResponse `json:"opportunities"`
	Links         []Link                `json:"links"`
	LastActiveAt  time.Time             `json:"last_active_at"`
	Username      string                `json:"username"`
} // @Name UserProfileResponse

func ToUserProfile(user db.User) UserProfileResponse {
	name := ""
	if user.Name != nil {
		name = *user.Name
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
		Name:          name,
		AvatarURL:     avatarURL,
		Title:         title,
		Description:   description,
		Location:      location,
		IsFollowing:   user.IsFollowing,
		Badges:        ToBadgeResponseList(user.Badges),
		Opportunities: ToOpportunityResponseList(user.Opportunities, user.LanguageCode),
		Links:         ToLinkResponseList(user.Links),
		LastActiveAt:  user.LastActiveAt,
		Username:      user.Username,
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

func ToLinkResponseList(links []db.Link) []Link {
	linkResponses := make([]Link, len(links))
	for i, link := range links {
		linkResponses[i] = Link{
			URL:   link.URL,
			Label: link.Label,
			Type:  link.Type,
			Order: link.Order,
		}
	}
	return linkResponses
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
	ID                 string                `json:"id"`
	UserID             string                `json:"user_id"`
	Title              string                `json:"title"`
	Description        string                `json:"description"`
	IsPayable          bool                  `json:"is_payable"`
	Badges             []BadgeResponse       `json:"badges"`
	Opportunity        OpportunityResponse   `json:"opportunity"`
	Location           *CityResponse         `json:"location"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
	User               UserProfileResponse   `json:"user"`
	VerificationStatus db.VerificationStatus `json:"verification_status"`
} // @Name CollaborationResponse

func ToCollaborationResponse(collab db.Collaboration) CollaborationResponse {
	// Create the response without the location first
	response := CollaborationResponse{
		ID:                 collab.ID,
		UserID:             collab.UserID,
		Title:              collab.Title,
		Description:        collab.Description,
		IsPayable:          collab.IsPayable,
		Badges:             ToBadgeResponseList(collab.Badges),
		Opportunity:        ToOpportunityResponseList([]db.Opportunity{collab.Opportunity}, db.LanguageEN)[0],
		CreatedAt:          collab.CreatedAt,
		UpdatedAt:          collab.UpdatedAt,
		User:               ToUserProfile(collab.User),
		VerificationStatus: collab.VerificationStatus,
	}

	if collab.Location != nil {
		resp := ToCityResponse(*collab.Location)
		response.Location = &resp
	}

	return response
}
