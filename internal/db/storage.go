package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage struct {
	client *mongo.Client
	db     *mongo.Database
}

func IsNoRowsError(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

func IsAlreadyExistsError(err error) bool {
	var mongoErr mongo.WriteException
	if errors.As(err, &mongoErr) {
		for _, writeErr := range mongoErr.WriteErrors {
			if writeErr.Code == 11000 { // MongoDB duplicate key error
				return true
			}
		}
	}
	return false
}

func ConnectDB(uri, db string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(25).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("failed to ping mongodb: %v", err)
	}

	if err := createIndexes(ctx, client.Database(db)); err != nil {
		log.Fatalf("failed to create indexes: %v", err)
	}

	return &Storage{
		client: client,
		db:     client.Database(db),
	}, nil
}

func (s *Storage) Database() *mongo.Database {
	return s.db
}

func (s *Storage) Client() *mongo.Client {
	return s.client
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type HealthStats struct {
	Status            string `json:"status"`
	Error             string `json:"error,omitempty"`
	Message           string `json:"message"`
	ActiveConnections int    `json:"active_connections"`
	MaxPoolSize       int    `json:"max_pool_size"`
}

func (s *Storage) Health() (HealthStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := HealthStats{}

	// Ping the database
	err := s.client.Ping(ctx, readpref.Primary())
	if err != nil {
		stats.Status = "down"
		stats.Error = fmt.Sprintf("db down: %v", err)
		return stats, fmt.Errorf("db down: %w", err)
	}

	stats.Status = "up"
	stats.Message = "It's healthy"

	return stats, nil
}
