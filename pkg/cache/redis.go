package cache

import (
	"context"
	"fmt"

	"github.com/shortlink/shortlink-service/internal/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// Init 初始化 Redis 连接
func Init(cfg *config.RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// 测试连接
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("Redis connected successfully")
	return nil
}

// Close 关闭 Redis 连接
func Close() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// GetClient 获取 Redis 客户端实例
func GetClient() *redis.Client {
	return RedisClient
}
