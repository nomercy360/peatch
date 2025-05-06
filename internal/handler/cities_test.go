package handler_test

import (
	"context"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func insertTestLocations(t *testing.T) {
	ctx := context.Background()
	coll := testutils.GetTestDBStorage().Database().Collection("cities")
	if _, err := coll.DeleteMany(ctx, bson.M{}); err != nil {
		t.Fatalf("Failed to clear cities collection: %v", err)
	}

	docs := []interface{}{
		bson.M{"_id": "1", "name": "CityA", "country_code": "CA", "country_name": "CountryA", "geo": bson.M{"type": "Point", "coordinates": []float64{-74.0060, 40.7128}}},
		bson.M{"_id": "2", "name": "CityB", "country_code": "CB", "country_name": "CountryB", "geo": bson.M{"type": "Point", "coordinates": []float64{20, 30}}},
		bson.M{"_id": "3", "name": "CityC", "country_code": "CC", "country_name": "CountryC", "geo": bson.M{"type": "Point", "coordinates": []float64{50, 60}}},
	}

	if _, err := coll.InsertMany(ctx, docs); err != nil {
		t.Fatalf("Failed to insert sample cities: %v", err)
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
