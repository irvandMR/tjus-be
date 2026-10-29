package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDatabaseConfig(env *Environment) *pgxpool.Pool {
	databaseURL := env.DatabaseUrl

	if databaseURL == "" {
		log.Fatalf("Database URL is empty")
	}
	poolConfig, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Fatalf("Gagal parse config database: %v", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Gagal inisialisasi pgxpool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Gagal ping database: %v", err)
	}

	fmt.Println("Database Connected successfully via pgxpool...")
	return pool

}
