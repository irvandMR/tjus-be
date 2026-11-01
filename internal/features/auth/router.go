package auth

import (
	"tjus-be/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetAuthRouter(rg *gin.RouterGroup, db *pgxpool.Pool, env *config.Environment, rdb *redis.Client) {
	authRepository := NewAuthRepository(db)
	authService := NewAuthService(authRepository, env, rdb)
	authHandler := NewAuthHandler(authService)
	authHandler.RegisterRoutes(rg.Group("/auth"))
}
