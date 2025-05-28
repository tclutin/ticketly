package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Storage interface {
	Set(ctx context.Context, ticketId, chatId uint64) error
	Get(ctx context.Context, ticketId uint64) (uint64, error)
	Delete(ctx context.Context, ticketID string) error
}

type RedisStorage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

func (r *RedisStorage) Set(ctx context.Context, ticketId, chatId uint64) error {
	key := fmt.Sprintf("chat:ticket:%d", ticketId)

	if err := r.client.Set(ctx, key, chatId, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisStorage) Get(ctx context.Context, ticketId uint64) (uint64, error) {
	key := fmt.Sprintf("chat:ticket:%d", ticketId)

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}

	chatId, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid chatID in redis: %w", err)
	}

	return chatId, nil
}

func (r *RedisStorage) Delete(ctx context.Context, ticketID string) error {
	return r.client.Del(ctx, ticketID).Err()
}
