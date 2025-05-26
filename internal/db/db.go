package db

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/mattn/go-sqlite3"
	"time"
)

//go:embed schema.sql
var schemaSQL string

type Storage struct {
	db *sql.DB
}

func init() {
	// Initialize sqlite-vec extension
	sqlite_vec.Auto()

	// Registers the sqlite3 driver with a ConnectHook so that we can
	// initialize the default PRAGMAs and load the sqlite-vec extension.
	//
	// Note 1: we don't define the PRAGMA as part of the dsn string
	// because not all pragmas are available.
	//
	// Note 2: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection.
	sql.Register("peatch_sqlite3",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				// Load sqlite-vec extension
				if err := conn.LoadExtension("vec0", ""); err != nil {
					// Extension might be already loaded, which is fine
					// Only return error if it's not a "already loaded" error
				}

				_, err := conn.Exec(`
					PRAGMA busy_timeout       = 10000;
					PRAGMA journal_mode       = WAL;
					PRAGMA journal_size_limit = 200000000;
					PRAGMA synchronous        = NORMAL;
					PRAGMA foreign_keys       = ON;
					PRAGMA temp_store         = MEMORY;
					PRAGMA cache_size         = -16000;
				`, nil)

				return err
			},
		},
	)
}

func NewStorage(dbFile string) (*Storage, error) {
	db, err := sql.Open("peatch_sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Storage{db: db}

	// Initialize schema
	if err := s.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return s, nil
}

func (s *Storage) initSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, schemaSQL)
	return err
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func (s *Storage) Close() error {
	return s.db.Close()
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type HealthStats struct {
	Status            string `json:"status"`
	Error             string `json:"error,omitempty"`
	Message           string `json:"message"`
	OpenConnections   int    `json:"open_connections"`
	InUse             int    `json:"in_use"`
	Idle              int    `json:"idle"`
	WaitCount         int64  `json:"wait_count"`
	WaitDuration      string `json:"wait_duration"`
	MaxIdleClosed     int64  `json:"max_idle_closed"`
	MaxLifetimeClosed int64  `json:"max_lifetime_closed"`
}

func (s *Storage) Health() (HealthStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := HealthStats{}

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats.Status = "down"
		stats.Error = fmt.Sprintf("db down: %v", err)
		return stats, fmt.Errorf("db down: %w", err)
	}

	// Database is up, add more statistics
	stats.Status = "up"
	stats.Message = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats.OpenConnections = dbStats.OpenConnections
	stats.InUse = dbStats.InUse
	stats.Idle = dbStats.Idle
	stats.WaitCount = dbStats.WaitCount
	stats.WaitDuration = dbStats.WaitDuration.String()
	stats.MaxIdleClosed = dbStats.MaxIdleClosed
	stats.MaxLifetimeClosed = dbStats.MaxLifetimeClosed

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats.Message = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats.Message = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats.Message = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats.Message = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats, nil
}
