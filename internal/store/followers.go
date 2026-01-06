package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type FollowerStore struct {
	db *sql.DB
}

type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"folloer_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *FollowerStore) Follow(ctx context.Context, followedID int64, userID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followedID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrorConflict
		}
	}
	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, followedID int64, userID int64) error {
	query := `DELETE FROM followers (user_id, follower_id)
	          WHERE user_id=$1 
			  AND follower_id=$2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followedID)

	if err != nil {
		return err
	}
	return nil
}
