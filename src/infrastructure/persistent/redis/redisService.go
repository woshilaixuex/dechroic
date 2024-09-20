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

// SetToMapString 将指定的键值对序列化存储
func (rs *RedisService) SetToMapString(ctx context.Context, key string, hashmap map[interface{}]interface{}) error {
	setMap := make(map[string]interface{}, len(hashmap))

	for k, v := range hashmap {
		switch k := k.(type) {
		case string:
			setMap[k] = v
		case int:
			setMap[strconv.Itoa(k)] = v
		case int64:
			setMap[strconv.FormatInt(k, 10)] = v
		default:
			return fmt.Errorf("unsupported key type: %T", k)
		}
	}
	jsonString, err := json.Marshal(setMap)
	if err != nil {
		return err
	}
	err = rs.client.Set(ctx, key, jsonString, 0).Err()
	return err
}

// GetMapString 获取反序列化map
func (rs *RedisService) GetMapString(ctx context.Context, key string) *redis.StringCmd {
	return rs.client.Get(ctx, key)
}

// GetFromMapString 获取反序列化map fieId对应的值
func (rs *RedisService) GetFromMapString(ctx context.Context, key, field string) (interface{}, error) {
	jsonString, err := rs.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	var hashmap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &hashmap); err != nil {
		return "", err
	}
	return hashmap[field], nil
}

// SetToMap 将指定的键值对添加到哈希表中
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

// GetFromMap 获取指定 key fueId 对应的 Map
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

// SetValue 储存键值对，如果是其他类型会进行序列化
func (rs *RedisService) SetValue(ctx context.Context, key string, value interface{}) error {
	// 检查值的类型
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		// 将复杂类型序列化为 JSON
		jsonData, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data = string(jsonData)
	}

	// 将值存储到 Redis
	return rs.client.Set(ctx, key, data, 0).Err()
}

// // GetValue 从 Redis 中获取指定 key 的值
func (rs *RedisService) GetValue(ctx context.Context, key string) (interface{}, error) {
	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // key 不存在，返回 nil
	} else if err != nil {
		return "", err // 其他错误
	}
	// 默认返回 string 类型
	return val, nil
}
