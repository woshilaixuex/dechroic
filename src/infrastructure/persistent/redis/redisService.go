package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: redis存储orm
 * @Date: 2024-06-14 16:46
 */

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

// SetArray 将数组序列化为 JSON 并存储到 Redis
func (rs *RedisService) SetArray(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.client.Set(ctx, key, jsonValue, 0).Err()
}

// GetArray 从 Redis 获取值并反序列化为指定类型的数组
func (rs *RedisService) GetArray(ctx context.Context, key string, result interface{}) error {
	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil // key 不存在，返回空
	} else if err != nil {
		return err // 其他错误
	}
	if err := json.Unmarshal([]byte(val), result); err != nil {
		return err
	}
	return nil
}

// AddToMap 将指定的键值对添加到哈希表中
func (rs *RedisService) SetToMap(ctx context.Context, key string, hashmap map[string]interface{}) error {
	return rs.client.HSet(ctx, key, hashmap).Err()
}

// GetMap 获取指定 key 的 Map
func (rs *RedisService) GetMap(ctx context.Context, key string) *redis.MapStringStringCmd {
	return rs.client.HGetAll(ctx, key)
}

// SetList 将一个切片推送到 Redis 列表中
func (rs *RedisService) SetList(ctx context.Context, key string, values []interface{}) error {
	pipe := rs.client.Pipeline()
	for _, value := range values {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		pipe.RPush(ctx, key, jsonValue)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// GetValue 从 Redis 中获取指定 key 的值，并反序列化为目标类型
func (rs *RedisService) GetValue(ctx context.Context, key string) ([]interface{}, error) {
	var result []any
	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return result, nil // key 不存在，返回空切片
	} else if err != nil {
		return nil, err // 其他错误
	}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}
