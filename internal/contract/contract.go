package contract

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peatch-io/peatch/internal/db"
)

type OkResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type StatusResponse struct {
	Success bool `json:"success"`
}

type AuthTelegramRequest struct {
	Query string `json:"query"`
}

func (a AuthTelegramRequest) Validate() error {
	if a.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	return nil
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  db.User `json:"user"`
}

type AuthAdminResponse struct {
	Token string `json:"token"`
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
