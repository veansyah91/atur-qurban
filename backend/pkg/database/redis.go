package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/username/qurban-app/config"
)

// ConnectRedis membuat koneksi ke Redis dan memverifikasi dengan Ping
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	// Verifikasi koneksi saat inisialisasi
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("ConnectRedis: ping gagal: %w", err)
	}

	return client, nil
}
