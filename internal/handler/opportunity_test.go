package handler_test

import (
	"context"
	"encoding/json"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/testutils"
	"net/http"
	"testing"
	"time"

	"github.com/peatch-io/peatch/internal/db"
)

func TestListOpportunities_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, 927635965, "mkkksim", "Maksim")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	token := authResp.Token

	opps := []db.Opportunity{
		{
			ID:            "badge1",
			Text:          "Test Badge 1",
			TextRU:        "Тестовый Бейдж 1",
			Description:   "Test Description 1",
			DescriptionRU: "Тестовое описание 1",
			Icon:          "icon1",
			Color:         "ff0000",
			CreatedAt:     time.Now(),
		},
		{
			ID:            "badge2",
			Text:          "Test Badge 2",
			TextRU:        "Тестовый Бейдж 2",
			Description:   "Test Description 2",
			DescriptionRU: "Тестовое описание 2",
			Icon:          "icon2",
			Color:         "00ff00",
			CreatedAt:     time.Now(),
		},
	}

	for _, opp := range opps {
		if err := ts.Storage.CreateOpportunity(context.Background(), opp); err != nil {
			t.Fatalf("failed to insert opportunity: %v", err)
		}
	}

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/opportunities", "", token, http.StatusOK)
	var respOpps []contract.OpportunityResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &respOpps); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respOpps) != 2 {
		t.Errorf("expected 2 opportunities, got %d", len(respOpps))
	}

	// Check that both opportunities exist (order may vary)
	foundOpp1 := false
	foundOpp2 := false
	for _, opp := range respOpps {
		if opp.Text == "Тестовый Бейдж 1" && opp.Description == "Тестовое описание 1" {
			foundOpp1 = true
		}
		if opp.Text == "Тестовый Бейдж 2" && opp.Description == "Тестовое описание 2" {
			foundOpp2 = true
		}
	}

	if !foundOpp1 {
		t.Error("Opportunity 'Тестовый Бейдж 1' not found in response")
	}

	if !foundOpp2 {
		t.Error("Opportunity 'Тестовый Бейдж 2' not found in response")
	}
}
