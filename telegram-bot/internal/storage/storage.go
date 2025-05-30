package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tclutin/ticketly/telegram_bot/internal/models"
)

type Storage interface {
	SetChatMeta(ctx context.Context, chatId int64, meta models.RealtimeChatMeta) error
	GetChatMeta(ctx context.Context, chatId int64) (models.RealtimeChatMeta, error)
	DeleteChatMeta(ctx context.Context, chatId int64) error
}

type RedisStorage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

func (r *RedisStorage) SetChatMeta(ctx context.Context, chatId int64, meta models.RealtimeChatMeta) error {
	key := fmt.Sprintf("chat:%d", chatId)

	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, 0).Err()
}

func (r *RedisStorage) GetChatMeta(ctx context.Context, chatId int64) (models.RealtimeChatMeta, error) {
	key := fmt.Sprintf("chat:%d", chatId)

	var meta models.RealtimeChatMeta

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return meta, nil
		}

		return models.RealtimeChatMeta{}, err
	}

	if err := json.Unmarshal([]byte(val), &meta); err != nil {
		return meta, err
	}

	return meta, nil
}

func (r *RedisStorage) DeleteChatMeta(ctx context.Context, chatId int64) error {
	key := fmt.Sprintf("chat:%d", chatId)
	return r.client.Del(ctx, key).Err()
}
