package auth

import (
	"time"
	"tjus-be/internal/config"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetAuthRouter(rg *gin.RouterGroup, db *pgxpool.Pool, env *config.Environment, rdb *redis.Client) {
	authRepository := NewAuthRepository(db)
	authService := NewAuthService(authRepository, env, rdb)
	authHandler := NewAuthHandler(authService)

	// Maksimal 5 percobaan login per menit per IP untuk mencegah Brute-Force
	loginLimiter := middleware.RateLimitMiddleware(rdb, "auth_login", 5, time.Minute)
	authHandler.RegisterRoutes(rg.Group("/auth"), loginLimiter)
}
