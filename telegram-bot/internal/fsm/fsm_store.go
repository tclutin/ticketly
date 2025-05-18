package fsm

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisFSMStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisFSMStore {
	return &RedisFSMStore{client: client}
}

func (s *RedisFSMStore) Get(userID int64) (string, error) {
	state, err := s.client.Get(context.Background(), fmt.Sprintf("fsm:%d", userID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return state, nil
}

func (s *RedisFSMStore) Set(userID int64, state string) error {
	return s.client.Set(context.Background(), fmt.Sprintf("fsm:%d", userID), state, 0).Err()
}

func (s *RedisFSMStore) Clear(userID int64) error {
	return s.client.Del(context.Background(), fmt.Sprintf("fsm:%d", userID)).Err()
}
