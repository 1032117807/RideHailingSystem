package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CodeRepository interface {
	Set(ctx context.Context, target, scene, code string, ttl time.Duration) error
	Verify(ctx context.Context, target, scene, code string) (bool, error)
}

type RedisCodeRepository struct {
	client redis.UniversalClient
	prefix string
}

func NewRedisCodeRepository(client redis.UniversalClient, prefix string) CodeRepository {
	return &RedisCodeRepository{
		client: client,
		prefix: prefix,
	}
}

func (r *RedisCodeRepository) Set(ctx context.Context, target, scene, code string, ttl time.Duration) error {
	key := r.buildKey(target, scene)
	return r.client.Set(ctx, key, code, ttl).Err()
}

func (r *RedisCodeRepository) Verify(ctx context.Context, target, scene, code string) (bool, error) {
	key := r.buildKey(target, scene)
	storedCode, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if storedCode != code {
		return false, nil
	}
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisCodeRepository) buildKey(target, scene string) string {
	return fmt.Sprintf("%s%s:%s", r.prefix, scene, strings.ToLower(strings.TrimSpace(target)))
}
