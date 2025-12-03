package redis

import (
	"context"
	"fmt"
	"merchant_api/pkg/config"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接 Redis 失败: %w", err)
	}

	return nil
}

// GetRedis 获取 Redis 客户端
func GetRedis() *redis.Client {
	return Client
}
