package cache

import (
	"FileNest/common/glog"
	"FileNest/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis 初始化 Redis 客户端
func InitRedis() error {
	cfg := config.Redis
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	glog.Infof("正在连接 Redis 服务器: %s", addr)

	redisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			glog.Info("Redis 连接成功")
			return nil
		},
	})

	// 测试连接并重试
	var lastErr error
	for i := 0; i <= cfg.MaxRetries; i++ {
		if i > 0 {
			glog.Infof("第 %d 次重试连接 Redis...", i)
			time.Sleep(time.Second * time.Duration(i)) // 重试间隔递增
		}

		if err := redisClient.Ping(ctx).Err(); err != nil {
			lastErr = fmt.Errorf("Redis 连接失败: %v", err)
			glog.Errorf("Redis 连接失败: %v, 地址: %s", err, addr)
			continue
		}

		glog.Infof("Redis 连接成功: %s, DB: %d", addr, cfg.DB)
		return nil
	}

	return lastErr
}

// GetRedisClient 获取 Redis 客户端
func GetRedisClient() *redis.Client {
	return redisClient
}

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	if err := redisClient.Set(ctx, key, value, expiration).Err(); err != nil {
		glog.Errorf("设置缓存失败: key=%s, error=%v", key, err)
		return err
	}
	return nil
}

// Get 获取缓存
func Get(key string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		glog.Errorf("获取缓存失败: key=%s, error=%v", key, err)
	}
	return val, err
}

// Del 删除缓存
func Del(keys ...string) error {
	if err := redisClient.Del(ctx, keys...).Err(); err != nil {
		glog.Errorf("删除缓存失败: keys=%v, error=%v", keys, err)
		return err
	}
	return nil
}

// DelByPattern 根据模式删除缓存
func DelByPattern(pattern string) error {
	keys, err := redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		glog.Errorf("查找缓存键失败: pattern=%s, error=%v", pattern, err)
		return err
	}
	if len(keys) > 0 {
		if err := redisClient.Del(ctx, keys...).Err(); err != nil {
			glog.Errorf("删除缓存失败: keys=%v, error=%v", keys, err)
			return err
		}
		glog.Infof("已删除 %d 个缓存键: pattern=%s", len(keys), pattern)
	}
	return nil
}

// HSet 设置哈希表字段
func HSet(key string, values ...interface{}) error {
	if err := redisClient.HSet(ctx, key, values...).Err(); err != nil {
		glog.Errorf("设置哈希表字段失败: key=%s, error=%v", key, err)
		return err
	}
	return nil
}

// HGet 获取哈希表字段
func HGet(key, field string) (string, error) {
	val, err := redisClient.HGet(ctx, key, field).Result()
	if err != nil && err != redis.Nil {
		glog.Errorf("获取哈希表字段失败: key=%s, field=%s, error=%v", key, field, err)
	}
	return val, err
}

// HGetAll 获取哈希表所有字段
func HGetAll(key string) (map[string]string, error) {
	val, err := redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		glog.Errorf("获取哈希表所有字段失败: key=%s, error=%v", key, err)
		return nil, err
	}
	return val, nil
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	if err := redisClient.Expire(ctx, key, expiration).Err(); err != nil {
		glog.Errorf("设置过期时间失败: key=%s, expiration=%v, error=%v", key, expiration, err)
		return err
	}
	return nil
}

// Close 关闭 Redis 连接
func Close() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
