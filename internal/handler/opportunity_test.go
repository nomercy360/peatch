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
	e := testutils.SetupHandlerDependencies(t)

	// Clear test data first
	if err := testutils.ClearTestData(); err != nil {
		t.Fatalf("failed to clear test data: %v", err)
	}

	authResp, err := testutils.AuthHelper(t, e, 927635965, "mkkksim", "Maksim")
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
		if err := testutils.GetDBStorage().CreateOpportunity(context.Background(), opp); err != nil {
			t.Fatalf("failed to insert opportunity: %v", err)
		}
	}

	rec := testutils.PerformRequest(t, e, http.MethodGet, "/api/opportunities", "", token, http.StatusOK)
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
		if opp.Text == "Тестовая Бейдж 1" && opp.Description == "Тестовое описание 1" {
			foundOpp1 = true
		}
		if opp.Text == "Тестовая Бейдж 2" && opp.Description == "Test Description 2" {
			foundOpp2 = true
		}
	}

	if !foundOpp1 {
		t.Error("expected to find opportunity with text 'Тестовая Бейдж 1' and description 'Тестовое описание 1'")
	}
	if !foundOpp2 {
		t.Error("expected to find opportunity with text 'Тестовая Бейдж 2' and description 'Test Description 2'")
	}
}
