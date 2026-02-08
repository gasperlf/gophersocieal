package cache

import (
	"context"

	"ontopsolutions.net/gasperlf/social/internal/store"
)

type MockCacheStore struct{}

func NewMockCache() Storage {
	return Storage{
		Users: &MockCacheStore{},
	}
}

func (m *MockCacheStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	return nil, nil
}

func (m *MockCacheStore) Set(ctx context.Context, user *store.User) error {
	return nil
}

func (m *MockCacheStore) Delete(ctx context.Context, userID int64) error {
	return nil
}
