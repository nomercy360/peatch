package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserFollower struct {
	ID         string    `bson:"_id,omitempty"`
	UserID     string    `bson:"user_id"`
	FollowerID string    `bson:"follower_id"`
	ExpiresAt  time.Time `bson:"expires_at"`
}
type LanguageCode string // @Name LanguageCode

var (
	LanguageEN LanguageCode = "en"
	LanguageRU LanguageCode = "ru"
)

type LoginMeta struct {
	IP        string `bson:"ip,omitempty" json:"ip,omitempty"`
	UserAgent string `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	Country   string `bson:"country,omitempty" json:"country,omitempty"`
	City      string `bson:"city,omitempty" json:"city,omitempty"`
}

type VerificationStatus string // @Name VerificationStatus

const (
	VerificationStatusPending    VerificationStatus = "pending"
	VerificationStatusVerified   VerificationStatus = "verified"
	VerificationStatusDenied     VerificationStatus = "denied"
	VerificationStatusBlocked    VerificationStatus = "blocked"
	VerificationStatusUnverified VerificationStatus = "unverified"
)

type User struct {
	ID                     string             `bson:"_id,omitempty" json:"id"`
	FirstName              *string            `bson:"first_name,omitempty" json:"first_name"`
	LastName               *string            `bson:"last_name,omitempty" json:"last_name"`
	ChatID                 int64              `bson:"chat_id,omitempty" json:"chat_id"`
	Username               string             `bson:"username,omitempty" json:"username"`
	CreatedAt              time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt              time.Time          `bson:"updated_at,omitempty" json:"-"`
	NotificationsEnabledAt *time.Time         `bson:"notifications_enabled_at,omitempty" json:"-"`
	HiddenAt               *time.Time         `bson:"hidden_at,omitempty" json:"hidden_at"`
	AvatarURL              *string            `bson:"avatar_url,omitempty" json:"avatar_url"`
	Title                  *string            `bson:"title,omitempty" json:"title"`
	Description            *string            `bson:"description,omitempty" json:"description"`
	LanguageCode           LanguageCode       `bson:"language_code,omitempty" json:"language_code"`
	Location               *City              `bson:"location,omitempty" json:"location"`
	IsFollowing            bool               `bson:"is_following,omitempty" json:"is_following"`
	Badges                 []Badge            `bson:"badges,omitempty" json:"badges"`
	Opportunities          []Opportunity      `bson:"opportunities,omitempty" json:"opportunities"`
	LoginMeta              *LoginMeta         `bson:"login_meta,omitempty" json:"login_meta"`
	LastActiveAt           time.Time          `bson:"last_active_at,omitempty" json:"last_active_at"`
	VerificationStatus     VerificationStatus `bson:"verification_status,omitempty" json:"verification_status"`
	VerifiedAt             *time.Time         `bson:"verified_at,omitempty" json:"verified_at"`
}

func (u *User) IsGeneratedUsername() bool {
	if strings.HasPrefix(u.Username, "user_") {
		return true
	}

	return false
}

func (u *User) IsProfileComplete() bool {
	if u.FirstName == nil || u.LastName == nil || u.Title == nil || u.Description == nil {
		return false
	}
	if u.Location == nil || u.Location.ID == "" {
		return false
	}
	if u.Badges == nil || len(u.Badges) == 0 {
		return false
	}
	if u.Opportunities == nil || len(u.Opportunities) == 0 {
		return false
	}
	if u.AvatarURL == nil || *u.AvatarURL == "" {
		return false
	}
	return true
}

type UserQuery struct {
	Page   int
	Limit  int
	Search string
	UserID string
}
type GetUsersParams struct {
	ViewerID string
	Username string
	UserID   string
	Lang     string
}

func (s *Storage) ListUsers(ctx context.Context, params UserQuery) ([]User, error) {
	collection := s.db.Collection("users")
	users := make([]User, 0)

	conditions := []bson.M{
		{"hidden_at": nil, "verification_status": VerificationStatusVerified},
	}

	if params.UserID != "" {
		conditions = append(conditions, bson.M{"_id": bson.M{"$ne": params.UserID}})
	}

	if params.Search != "" {
		searchRegex := primitive.Regex{Pattern: params.Search, Options: "i"}
		conditions = append(conditions, bson.M{
			"$or": []bson.M{
				{"first_name": searchRegex},
				{"last_name": searchRegex},
				{"title": searchRegex},
				{"description": searchRegex},
			},
		})
	}

	var filter bson.M
	if len(conditions) == 1 {
		filter = conditions[0]
	} else {
		filter = bson.M{"$and": conditions}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(params.Limit))
	findOptions.SetSkip(int64((params.Page - 1) * params.Limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

func (s *Storage) getUserBy(ctx context.Context, filter bson.M) (User, error) {
	collection := s.db.Collection("users")

	var user User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return user, nil
}

func (s *Storage) GetUserByChatID(ctx context.Context, chatID int64) (User, error) {
	return s.getUserBy(ctx, bson.M{"chat_id": chatID})
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (User, error) {
	return s.getUserBy(ctx, bson.M{"_id": id})
}

func (s *Storage) CreateUser(ctx context.Context, user User) error {
	collection := s.db.Collection("users")

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.LastActiveAt = now
	user.VerificationStatus = VerificationStatusUnverified

	if _, err := collection.InsertOne(ctx, user); err != nil {
		return nil
	}

	return nil
}

func (s *Storage) UpdateUser(ctx context.Context, user User, badgeIDs, oppIDs []string, locationID string) error {
	collection := s.db.Collection("users")
	badgeCollection := s.db.Collection("badges")
	oppCollection := s.db.Collection("opportunities")
	locationCollection := s.db.Collection("cities")

	badgeFilter := bson.M{"_id": bson.M{"$in": badgeIDs}}
	badgeCursor, err := badgeCollection.Find(ctx, badgeFilter)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNotFound
	} else if err != nil {
		return nil
	}

	var badges []Badge
	if err := badgeCursor.All(ctx, &badges); err != nil {
		return fmt.Errorf("failed to decode badges: %w", err)
	}

	var opps []Opportunity
	oppFilter := bson.M{"_id": bson.M{"$in": oppIDs}}
	oppCursor, err := oppCollection.Find(ctx, oppFilter)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNotFound
	} else if err != nil {
		return nil
	}

	if err := oppCursor.All(ctx, &opps); err != nil {
		return fmt.Errorf("failed to decode opportunities: %w", err)
	}

	var locationData City
	locationFilter := bson.M{"_id": locationID}
	if err := locationCollection.FindOne(ctx, locationFilter).Decode(&locationData); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}

		return fmt.Errorf("failed to fetch location: %w", err)
	}

	filter := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"updated_at":    time.Now(),
			"title":         user.Title,
			"description":   user.Description,
			"location":      locationData,
			"badges":        badges,
			"opportunities": opps,
		},
	}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Storage) GetUsersByVerificationStatus(ctx context.Context, status VerificationStatus, page, perPage int) ([]User, error) {
	collection := s.db.Collection("users")

	filter := bson.M{"verification_status": status}

	skip := (page - 1) * perPage

	findOptions := options.Find().
		SetLimit(int64(perPage)).
		SetSkip(int64(skip)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by verification status: %w", err)
	}
	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

func (s *Storage) GetUserProfile(ctx context.Context, viewerID string, id string) (User, error) {
	usersCollection := s.db.Collection("users")
	followersCollection := s.db.Collection("user_followers")

	var user User

	filter := bson.M{
		"_id": id,
		"$or": []bson.M{
			{"_id": viewerID},
			{
				"$and": []bson.M{
					{"verification_status": VerificationStatusVerified},
					{"hidden_at": nil},
				},
			},
		},
	}

	err := usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, ErrNotFound
		}
		return user, fmt.Errorf("failed to get user profile: %w", err)
	}

	followerFilter := bson.M{
		"user_id":     user.ID,
		"follower_id": viewerID,
	}

	count, err := followersCollection.CountDocuments(ctx, followerFilter)
	if err != nil {
		user.IsFollowing = false
	} else {
		user.IsFollowing = count > 0
	}

	return user, nil
}

func (s *Storage) FollowUser(ctx context.Context, userID string, followerID string, ttlDuration time.Duration) error {
	usersCollection := s.db.Collection("users")
	followersCollection := s.db.Collection("user_followers")

	userFilter := bson.M{
		"_id":                 userID,
		"hidden_at":           nil,
		"verification_status": VerificationStatusVerified,
	}

	count, err := usersCollection.CountDocuments(ctx, userFilter)
	if err != nil {
		return fmt.Errorf("failed to check user existence for follow: %w", err)
	}

	if count == 0 {
		return ErrNotFound
	}

	expiresAt := time.Now().Add(ttlDuration)

	followDoc := UserFollower{
		UserID:     userID,
		FollowerID: followerID,
		ExpiresAt:  expiresAt,
	}

	_, err = followersCollection.InsertOne(ctx, followDoc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return fmt.Errorf("failed to follow user: %w", err)
	}

	return nil
}

func (s *Storage) IsUserFollowing(ctx context.Context, userID string, followerID string) (bool, error) {
	followersCollection := s.db.Collection("user_followers")

	filter := bson.M{
		"follower_id": followerID,
		"user_id":     userID,
		"expires_at":  bson.M{"$gt": time.Now()},
	}

	count, err := followersCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check following status: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) UpdateUserAvatarURL(ctx context.Context, userID string, avatarURL string) error {
	collection := s.db.Collection("users")

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"avatar_url": avatarURL}}

	res, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("failed to update user avatar URL: %w", err)
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Storage) UpdateUserLoginMetadata(ctx context.Context, userID string, meta LoginMeta) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := s.db.Collection("users")
	now := time.Now()

	update := bson.M{
		"last_login":     meta,
		"last_active_at": now,
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": update},
	)
	return err
}

func (s *Storage) UpdateUserVerificationStatus(ctx context.Context, userID string, status VerificationStatus) error {
	collection := s.db.Collection("users")
	now := time.Now()

	update := bson.M{
		"$set": bson.M{
			"verification_status": status,
			"updated_at":          now,
		},
	}

	if status == VerificationStatusVerified {
		update = bson.M{
			"$set": bson.M{
				"verification_status": status,
				"verified_at":         time.Now(),
				"updated_at":          time.Now(),
			},
		}
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}
