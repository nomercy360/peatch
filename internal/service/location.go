package service

import "github.com/peatch-io/peatch/internal/db"

func (s *service) SearchLocations(query string) ([]db.Location, error) {
	res, err := s.storage.SearchLocations(query)

	if err != nil {
		return nil, err
	}

	return res, nil
}
