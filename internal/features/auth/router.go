package auth

import (
	"tjus-be/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetAuthRouter(rg *gin.RouterGroup, db *pgxpool.Pool, env *config.Environment) {
	authRepository := NewAuthRepository(db)
	authService := NewAuthService(authRepository, env)
	authHandler := NewAuthHandler(authService)
	authHandler.RegisterRoutes(rg.Group("/auth"))
}
