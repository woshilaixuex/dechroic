package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	client    *redis.Client
	BatchSize int
}

func NewRedisService(addr, password string, db int) *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,               // 0~15默认为0
		DialTimeout:  10 * time.Second, // 连接超时
		ReadTimeout:  10 * time.Second, // 读取超时
		WriteTimeout: 10 * time.Second, // 写入超时
		PoolSize:     10,               // 连接池大小
		MinIdleConns: 5,                // 最小空闲连接数
		MaxRetries:   3,                // 重试次数
	})
	return &RedisService{
		client:    rdb,
		BatchSize: 1000,
	}
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
func (rs *RedisService) SetToMap(ctx context.Context, key string, hashmap map[interface{}]interface{}) error {
	setMap := make(map[string]interface{}, len(hashmap))

	for k, v := range hashmap {
		switch k := k.(type) {
		case string:
			setMap[k] = v
		case int:
			setMap[strconv.Itoa(k)] = v
		case int64:
			setMap[strconv.FormatInt(k, 10)] = v
		// 添加更多类型转换根据需要
		default:
			return fmt.Errorf("unsupported key type: %T", k)
		}
	}
	// 分批次存储数据
	batch := make(map[string]interface{}, rs.BatchSize)
	count := 0
	for k, v := range setMap {
		batch[k] = v
		count++
		if count == rs.BatchSize {
			if err := rs.client.HSet(ctx, key, batch).Err(); err != nil {
				return err
			}
			batch = make(map[string]interface{}, rs.BatchSize)
			count = 0
		}
	}
	// 存储剩余数据
	if len(batch) > 0 {
		if err := rs.client.HSet(ctx, key, batch).Err(); err != nil {
			return err
		}
	}
	return nil
}

// GetMap 获取指定 key 的 Map
func (rs *RedisService) GetMap(ctx context.Context, key string) *redis.MapStringStringCmd {
	return rs.client.HGetAll(ctx, key)
}
func (rs *RedisService) GetFromMap(ctx context.Context, key, field string) *redis.StringCmd {
	return rs.client.HGet(ctx, key, field)
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
func (rs *RedisService) GetValue(ctx context.Context, key string) (interface{}, error) {
	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // key 不存在，返回 nil
	} else if err != nil {
		return nil, err // 其他错误
	}

	// 尝试将返回值转换为 int64 或 int
	if int64Val, err := strconv.ParseInt(val, 10, 64); err == nil {
		return int64Val, nil
	}
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal, nil
	}

	// 默认返回 string 类型
	return val, nil
}
