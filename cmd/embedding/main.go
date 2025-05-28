package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/peatch-io/peatch/internal/config"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/embedding"
)

func main() {
	var (
		dbPath     = flag.String("db", "data.sqlite", "Path to SQLite database")
		dryRun     = flag.Bool("dry-run", false, "Only show what would be updated")
		userOnly   = flag.Bool("users-only", false, "Only generate user embeddings")
		collabOnly = flag.Bool("collab-only", false, "Only generate collaboration embeddings")
		limit      = flag.Int("limit", 0, "Limit number of records to process (0 = no limit)")
		batchSize  = flag.Int("batch", 10, "Number of embeddings to generate in parallel")
		sleep      = flag.Duration("sleep", 100*time.Millisecond, "Sleep duration between batches")
	)
	flag.Parse()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.OpenAIAPIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create storage
	storage, err := db.NewStorage(*dbPath)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.DB().Close()

	// Create embedding service
	embeddingService := embedding.New(cfg.OpenAIAPIKey)

	ctx := context.Background()

	// Process based on flags
	if !*collabOnly {
		if err := processUsers(ctx, storage, embeddingService, *dryRun, *limit, *batchSize, *sleep); err != nil {
			log.Fatalf("Failed to process users: %v", err)
		}
	}

	if !*userOnly {
		if err := processCollaborations(ctx, storage, embeddingService, *dryRun, *limit, *batchSize, *sleep); err != nil {
			log.Fatalf("Failed to process collaborations: %v", err)
		}
	}
}

func processUsers(ctx context.Context, storage *db.Storage, embSvc *embedding.Service, dryRun bool, limit, batchSize int, sleep time.Duration) error {
	log.Println("Processing users without embeddings...")

	// Get users needing embeddings
	users, err := getUsersWithoutEmbeddings(ctx, storage, limit)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	log.Printf("Found %d users without embeddings", len(users))

	if dryRun {
		for _, user := range users {
			log.Printf("Would generate embedding for user: %s (%s)", *user.Name, user.ID)
		}
		return nil
	}

	// Process in batches
	processed := 0
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		for _, user := range batch {
			// Generate embedding
			vec, err := embSvc.GenerateEmbedding(ctx, user.ToString())
			if err != nil {
				log.Printf("Failed to generate embedding for user %s: %v", user.ID, err)
				continue
			}

			// Store embedding
			if err := storage.UpdateUserEmbedding(ctx, user.ID, vec); err != nil {
				log.Printf("Failed to store embedding for user %s: %v", user.ID, err)
				continue
			}

			processed++
			log.Printf("Generated embedding for user: %s (%s) [%d/%d]", *user.Name, user.ID, processed, len(users))
		}

		// Sleep between batches to avoid rate limits
		if i+batchSize < len(users) {
			time.Sleep(sleep)
		}
	}

	log.Printf("Successfully processed %d/%d users", processed, len(users))
	return nil
}

func processCollaborations(ctx context.Context, storage *db.Storage, embSvc *embedding.Service, dryRun bool, limit, batchSize int, sleep time.Duration) error {
	log.Println("Processing collaborations...")

	// Get collaborations without embeddings
	collabs, err := getCollaborationsWithoutEmbeddings(ctx, storage, limit)
	if err != nil {
		return fmt.Errorf("failed to get collaborations: %w", err)
	}

	log.Printf("Found %d collaborations without embeddings", len(collabs))

	if dryRun {
		for _, collab := range collabs {
			log.Printf("Would generate embedding for collaboration: %s (%s)", collab.Title, collab.ID)
		}
		return nil
	}

	// Process in batches
	processed := 0
	for i := 0; i < len(collabs); i += batchSize {
		end := i + batchSize
		if end > len(collabs) {
			end = len(collabs)
		}

		batch := collabs[i:end]
		for _, collab := range batch {
			// Build embedding text
			text := collab.ToString()
			if text == "" {
				log.Printf("Skipping collaboration %s - no content to embed", collab.ID)
				continue
			}

			// Generate embedding
			vec, err := embSvc.GenerateEmbedding(ctx, text)
			if err != nil {
				log.Printf("Failed to generate embedding for collaboration %s: %v", collab.ID, err)
				continue
			}

			// Store embedding
			if err := updateCollaborationEmbedding(ctx, storage, collab.ID, vec); err != nil {
				log.Printf("Failed to store embedding for collaboration %s: %v", collab.ID, err)
				continue
			}

			processed++
			log.Printf("Generated embedding for collaboration: %s (%s) [%d/%d]", collab.Title, collab.ID, processed, len(collabs))
		}

		// Sleep between batches
		if i+batchSize < len(collabs) {
			time.Sleep(sleep)
		}
	}

	log.Printf("Successfully processed %d/%d collaborations", processed, len(collabs))
	return nil
}

func getUsersWithoutEmbeddings(ctx context.Context, storage *db.Storage, limit int) ([]db.User, error) {
	query := `
		SELECT u.id, u.chat_id, u.name, u.username, u.avatar_url, 
		       u.title, u.description, u.location, u.badges, u.opportunities,
		       u.links, u.created_at, u.updated_at, u.hidden_at, 
		       u.verification_status, u.verified_at, u.embedding_updated_at
		FROM users u
		LEFT JOIN user_embeddings ue ON u.id = ue.user_id
		WHERE ue.user_id IS NULL
		  AND (u.name IS NOT NULL AND u.description IS NOT NULL AND u.title IS NOT NULL AND u.badges IS NOT NULL AND u.opportunities IS NOT NULL)
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := storage.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []db.User
	for rows.Next() {
		var user db.User
		var locationJSON, badgesJSON, oppsJSON, linksJSON sql.NullString

		err := rows.Scan(
			&user.ID, &user.ChatID, &user.Name, &user.Username, &user.AvatarURL,
			&user.Title, &user.Description, &locationJSON, &badgesJSON, &oppsJSON,
			&linksJSON, &user.CreatedAt, &user.UpdatedAt, &user.HiddenAt,
			&user.VerificationStatus, &user.VerifiedAt, &user.EmbeddingUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if locationJSON.Valid && locationJSON.String != "" {
			json.Unmarshal([]byte(locationJSON.String), &user.Location)
		}
		if badgesJSON.Valid && badgesJSON.String != "" {
			json.Unmarshal([]byte(badgesJSON.String), &user.Badges)
		}
		if oppsJSON.Valid && oppsJSON.String != "" {
			json.Unmarshal([]byte(oppsJSON.String), &user.Opportunities)
		}
		if linksJSON.Valid && linksJSON.String != "" {
			json.Unmarshal([]byte(linksJSON.String), &user.Links)
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

func getCollaborationsWithoutEmbeddings(ctx context.Context, storage *db.Storage, limit int) ([]db.Collaboration, error) {
	query := `
		SELECT c.id, c.user_id, c.title, c.description, c.is_payable,
		       c.created_at, c.updated_at, c.hidden_at, c.location, 
		       c.links, c.badges, c.opportunity, c.verification_status, c.verified_at
		FROM collaborations c
		LEFT JOIN collaboration_embeddings ce ON c.id = ce.collaboration_id
		WHERE ce.collaboration_id IS NULL
		  AND c.hidden_at IS NULL
		  AND c.verification_status = 'verified'
		  AND (c.title IS NOT NULL AND c.description IS NOT NULL)
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := storage.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collabs []db.Collaboration
	for rows.Next() {
		var collab db.Collaboration
		var locationJSON, linksJSON, badgesJSON, opportunityJSON sql.NullString

		err := rows.Scan(
			&collab.ID, &collab.UserID, &collab.Title, &collab.Description, &collab.IsPayable,
			&collab.CreatedAt, &collab.UpdatedAt, &collab.HiddenAt, &locationJSON,
			&linksJSON, &badgesJSON, &opportunityJSON, &collab.VerificationStatus, &collab.VerifiedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if locationJSON.Valid && locationJSON.String != "" {
			var location db.City
			if err := json.Unmarshal([]byte(locationJSON.String), &location); err == nil {
				collab.Location = &location
			}
		}
		if linksJSON.Valid && linksJSON.String != "" {
			json.Unmarshal([]byte(linksJSON.String), &collab.Links)
		}
		if badgesJSON.Valid && badgesJSON.String != "" {
			json.Unmarshal([]byte(badgesJSON.String), &collab.Badges)
		}
		if opportunityJSON.Valid && opportunityJSON.String != "" {
			json.Unmarshal([]byte(opportunityJSON.String), &collab.Opportunity)
		}

		collabs = append(collabs, collab)
	}

	return collabs, rows.Err()
}

func updateCollaborationEmbedding(ctx context.Context, storage *db.Storage, collabID string, vec []float64) error {
	tx, err := storage.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Serialize embedding for sqlite-vec
	serialized, err := sqlite_vec.SerializeFloat32(vectorFloat64ToFloat32(vec))
	if err != nil {
		return fmt.Errorf("failed to serialize vector: %w", err)
	}

	// insert embedding
	_, err = tx.ExecContext(ctx, `
		INSERT INTO collaboration_embeddings (collaboration_id, embedding)
		VALUES (?, ?)
	`, collabID, serialized)

	if err != nil {
		return fmt.Errorf("failed to update collaboration embedding: %w", err)
	}

	// Update timestamp
	_, err = tx.ExecContext(ctx, `
		UPDATE collaborations 
		SET embedding_updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, collabID)

	if err != nil {
		return fmt.Errorf("failed to update collaboration timestamp: %w", err)
	}

	return tx.Commit()
}

// vectorFloat64ToFloat32 converts []float64 to []float32 for sqlite-vec
func vectorFloat64ToFloat32(vec []float64) []float32 {
	result := make([]float32, len(vec))
	for i, v := range vec {
		result[i] = float32(v)
	}
	return result
}
