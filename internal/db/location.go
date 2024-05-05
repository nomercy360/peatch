package db

type Location struct {
	ID          int64  `json:"-" db:"id"`
	City        string `json:"city" db:"city"`
	Country     string `json:"country" db:"country"`
	CountryCode string `json:"country_code" db:"country_code"`
	Population  int64  `json:"-" db:"population"`
}

func (s *storage) SearchLocations(query string) ([]Location, error) {
	locations := make([]Location, 0)

	err := s.pg.Select(&locations, "SELECT id, city, country_code, country, population FROM locations WHERE city ILIKE $1 OR country ILIKE $1 ORDER BY population desc LIMIT 20", "%"+query+"%")

	return locations, err
}
