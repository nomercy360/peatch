package db

import (
	"context"
	"fmt"
	"time"
)

// CleanupExpiredRecords removes expired records from tables with TTL
func (s *Storage) CleanupExpiredRecords(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// Clean up expired user followers
	result, err := tx.ExecContext(ctx,
		`DELETE FROM user_followers WHERE expires_at < ?`, now)
	if err != nil {
		return fmt.Errorf("failed to cleanup user_followers: %w", err)
	}
	followersDeleted, _ := result.RowsAffected()

	// Clean up expired collaboration interests
	result, err = tx.ExecContext(ctx,
		`DELETE FROM collaboration_interests WHERE expires_at < ?`, now)
	if err != nil {
		return fmt.Errorf("failed to cleanup collaboration_interests: %w", err)
	}
	interestsDeleted, _ := result.RowsAffected()

	// Log cleanup results
	_, err = tx.ExecContext(ctx, `
		INSERT INTO cleanup_log (table_name, deleted) VALUES 
		('user_followers', ?),
		('collaboration_interests', ?)
	`, followersDeleted, interestsDeleted)
	if err != nil {
		return fmt.Errorf("failed to log cleanup: %w", err)
	}

	return tx.Commit()
}

// GetCleanupStats returns statistics about cleanup operations
func (s *Storage) GetCleanupStats(ctx context.Context, since time.Time) ([]CleanupStat, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT table_name, SUM(deleted) as total_deleted, COUNT(*) as cleanup_runs
		FROM cleanup_log
		WHERE executed_at >= ?
		GROUP BY table_name
	`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []CleanupStat
	for rows.Next() {
		var stat CleanupStat
		if err := rows.Scan(&stat.TableName, &stat.TotalDeleted, &stat.CleanupRuns); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

type CleanupStat struct {
	TableName    string
	TotalDeleted int64
	CleanupRuns  int
}
