package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type City struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
}

// SearchCities searches cities by name with pagination
func (s *Storage) SearchCities(ctx context.Context, search string, limit, skip int) ([]City, error) {
	query := `
		SELECT id, name, country_code, country_name, latitude, longitude
		FROM cities
		WHERE 1=1
	`
	var args []interface{}

	// Add search filter
	if search != "" {
		query += fmt.Sprintf(` AND name LIKE ?`)
		args = append(args, "%"+search+"%")
	}

	// Add ordering and pagination
	query += fmt.Sprintf(` ORDER BY name ASC LIMIT ? OFFSET ?`)
	args = append(args, limit, skip)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		city, err := scanCity(rows)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, rows.Err()
}

// GetCityByID retrieves a city by ID
func (s *Storage) GetCityByID(ctx context.Context, id string) (City, error) {
	query := `
		SELECT id, name, country_code, country_name, latitude, longitude
		FROM cities
		WHERE id = ?
	`

	row := s.db.QueryRowContext(ctx, query, id)
	city, err := scanCityRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return City{}, ErrNotFound
		}
		return City{}, err
	}

	return city, nil
}

// CreateCity creates a new city
func (s *Storage) CreateCity(ctx context.Context, city City) error {
	query := `
		INSERT INTO cities (id, name, country_code, country_name, latitude, longitude, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		city.ID,
		city.Name,
		city.CountryCode,
		city.CountryName,
		city.Latitude,
		city.Longitude,
		time.Now(),
	)

	if err != nil {
		if isSQLiteConstraintError(err) {
			return ErrAlreadyExists
		}
		return err
	}

	return nil
}

func scanCity(rows *sql.Rows) (City, error) {
	var city City

	err := rows.Scan(
		&city.ID,
		&city.Name,
		&city.CountryCode,
		&city.CountryName,
		&city.Latitude,
		&city.Longitude,
	)
	if err != nil {
		return city, err
	}

	return city, nil
}

func scanCityRow(row *sql.Row) (City, error) {
	var city City

	err := row.Scan(
		&city.ID,
		&city.Name,
		&city.CountryCode,
		&city.CountryName,
		&city.Latitude,
		&city.Longitude,
	)
	if err != nil {
		return city, err
	}

	return city, nil
}

func (s *Storage) fetchCityTx(ctx context.Context, tx *sql.Tx, id string) (City, error) {
	query := `
		SELECT id, name, country_code, country_name, latitude, longitude
		FROM cities
		WHERE id = ?
	`

	row := tx.QueryRowContext(ctx, query, id)
	city, err := scanCityRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return City{}, ErrNotFound
		}
		return City{}, err
	}

	return city, nil
}
