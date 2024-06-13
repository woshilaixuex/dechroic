package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(addr, password string, db int) *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db, // 0~15默认为0
	})
	return &RedisService{client: rdb}
}

// SetValue 设置指定 key 的值
func (rs *RedisService) SetValue(ctx context.Context, key string, value interface{}) error {
	return rs.client.Set(ctx, key, value, 0).Err()
}

// SetValueWithExpiration 设置指定 key 的值，并设置过期时间
func (rs *RedisService) SetValueWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rs.client.Set(ctx, key, value, expiration).Err()
}

// GetMap 获取指定 key 的 Map
func (rs *RedisService) GetMap(ctx context.Context, key string) *redis.MapStringStringCmd {
	return rs.client.HGetAll(ctx, key)
}

// AddToMap 将指定的键值对添加到哈希表中
func (rs *RedisService) AddToMap(ctx context.Context, key string, hashmap map[string]any) error {
	return rs.client.HSet(ctx, key, hashmap).Err()
}
