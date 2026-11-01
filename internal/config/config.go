package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Env   *Environment
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func InitConfig() *Config {
	env := LoadEnvironment()
	db := InitDatabaseConfig(env)
	RunMigrations(db)
	rdb := InitRedisConfig(env)

	return &Config{
		Env:   env,
		DB:    db,
		Redis: rdb,
	}
}

func (c *Config) Close() {
	if c.Redis != nil {
		_ = c.Redis.Close()
	}
	if c.DB != nil {
		c.DB.Close()
	}
}
