package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"ontopsolutions.net/gasperlf/social/internal/store"
)

const userExp = time.Minute

type UserStore struct {
	rdb *redis.Client
}

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err // Redis error
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEx(ctx, cacheKey, data, userExp).Err()
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	cacheKey := fmt.Sprintf("user-%v", userID)
	return s.rdb.Del(ctx, cacheKey).Err()
}
