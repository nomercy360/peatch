package handler_test

import (
	"context"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func insertTestLocations(t *testing.T) {
	ctx := context.Background()
	storage := testutils.GetDBStorage()

	// Clear existing cities
	if _, err := storage.DB().ExecContext(ctx, "DELETE FROM cities WHERE 1=1"); err != nil {
		t.Fatalf("Failed to clear cities table: %v", err)
	}

	cities := []db.City{
		{ID: "1", Name: "CityA", CountryCode: "CA", CountryName: "CountryA", Latitude: -74.0060, Longitude: 40.7128},
		{ID: "2", Name: "CityB", CountryCode: "CB", CountryName: "CountryB", Latitude: 30.0, Longitude: 40.0},
		{ID: "3", Name: "CityC", CountryCode: "CC", CountryName: "CountryC", Latitude: 50.0, Longitude: 60.0},
	}

	for _, city := range cities {
		if err := storage.CreateCity(ctx, city); err != nil {
			t.Fatalf("Failed to insert city %s: %v", city.ID, err)
		}
	}
}

func TestSearchCities(t *testing.T) {
	e := testutils.SetupHandlerDependencies(t)

	auth, err := testutils.AuthHelper(t, e, testutils.TelegramTestUserID, "tester", "Tester")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	insertTestLocations(t)

	rec := testutils.PerformRequest(t, e, http.MethodGet, "/api/locations?limit=2&page=2", "", auth.Token, http.StatusOK)
	resp := testutils.ParseResponse[[]contract.CityResponse](t, rec)
	assert.Len(t, resp, 1, "expected one city on second page with limit=2")
	assert.Equal(t, "3", resp[0].ID)
	assert.Equal(t, "CityC", resp[0].Name)

	rec = testutils.PerformRequest(t, e, http.MethodGet, "/api/locations?search=cityb", "", auth.Token, http.StatusOK)
	resp = testutils.ParseResponse[[]contract.CityResponse](t, rec)
	assert.Len(t, resp, 1, "expected one city matching 'cityb'")
	assert.Equal(t, "2", resp[0].ID)
	assert.Equal(t, "CityB", resp[0].Name)
}
