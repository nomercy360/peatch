package handler_test

import (
	"context"
	"encoding/json"
	"github.com/peatch-io/peatch/internal/contract"
	"net/http"
	"testing"
	"time"

	"github.com/peatch-io/peatch/internal/db"
)

func TestListOpportunities_Success(t *testing.T) {
	e := setupDependencies(t)

	authResp, err := authHelper(t, e, 927635965, "mkkksim", "Maksim")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	token := authResp.Token

	opps := []db.Opportunity{
		{
			ID:            "badge1",
			Text:          "Test Badge 1",
			TextRU:        "Тестовая Бейдж 1",
			Description:   "Test Description 1",
			DescriptionRU: "Тестовое описание 1",
			Icon:          "icon1",
			Color:         "ff0000",
			CreatedAt:     time.Now(),
		},
		{
			ID:          "badge2",
			Text:        "Test Badge 2",
			TextRU:      "Тестовая Бейдж 2",
			Description: "Test Description 2",
			Icon:        "icon2",
			Color:       "00ff00",
			CreatedAt:   time.Now(),
		},
	}

	for _, opp := range opps {
		if _, err := dbStorage.Database().Collection(OpportunitiesCollection).InsertOne(context.Background(), opp); err != nil {
			t.Fatalf("failed to insert opportunity: %v", err)
		}
	}

	rec := performRequest(t, e, http.MethodGet, "/api/opportunities", "", token, http.StatusOK)
	var respOpps []contract.OpportunityResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &respOpps); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respOpps) != 2 {
		t.Errorf("expected 2 badges, got %d", len(respOpps))
	}

	if respOpps[0].Text != "Тестовая Бейдж 1" {
		t.Errorf("expected badge text 'Test Badge 1', got '%s'", respOpps[0].Text)
	}

	if respOpps[0].Description != "Тестовое описание 1" {
		t.Errorf("expected badge text 'Test Description 1', got '%s'", respOpps[0].Description)
	}
}
