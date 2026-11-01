package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedisConfig(env *Environment) *redis.Client {
	if env.RedisAddr == "" {
		log.Println("REDIS_ADDR kosong, aplikasi berjalan tanpa Redis cache")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr: env.RedisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Peringatan: Gagal ping Redis (%v), aplikasi berjalan tanpa cache", err)
		_ = client.Close()
		return nil
	}
	fmt.Println("Redis Connected successfully...")
	return client
}
