package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string    `bson:"username" json:"username"`
	ChatID    int64     `bson:"chat_id" json:"chat_id"`
	Password  string    `bson:"password" json:"-"` // Never expose password in JSON responses
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (s *Storage) CreateAdmin(ctx context.Context, admin Admin) (Admin, error) {
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	hashedPassword, err := hashPassword(admin.Password)
	if err != nil {
		return Admin{}, err
	}

	admin.Password = hashedPassword

	collection := s.db.Collection("admins")

	if _, err := collection.InsertOne(ctx, admin); err != nil {
		return Admin{}, err
	}

	return admin, nil
}

func (s *Storage) GetAdminByUsername(ctx context.Context, username string) (Admin, error) {
	collection := s.db.Collection("admins")
	filter := bson.M{"username": username}

	var admin Admin
	err := collection.FindOne(ctx, filter).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Admin{}, errors.New("admin not found")
		}
		return Admin{}, err
	}

	return admin, nil
}

func (s *Storage) GetAdminByChatID(ctx context.Context, chatID int64) (Admin, error) {
	collection := s.db.Collection("admins")
	filter := bson.M{"chat_id": chatID}

	var admin Admin
	err := collection.FindOne(ctx, filter).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Admin{}, ErrNotFound
		}
		return Admin{}, err
	}

	return admin, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Storage) ValidateAdminCredentials(ctx context.Context, username, password string) (Admin, error) {
	admin, err := s.GetAdminByUsername(ctx, username)
	if err != nil {
		return Admin{}, err
	}

	if !checkPasswordHash(password, admin.Password) {
		return Admin{}, errors.New("invalid credentials")
	}

	return admin, nil
}
