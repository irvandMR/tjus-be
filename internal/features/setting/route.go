package setting

import (
	"tjus-be/internal/config"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetSettingRouting(rg *gin.RouterGroup, db *pgxpool.Pool, env *config.Environment, rdb *redis.Client) {
	repository := NewSettingsRepository(db)
	service := NewSettingsService(repository)
	handler := NewSettingsHandler(service)

	authMiddleware := middleware.AuthMiddleware(rdb, env)

	// Semua endpoint setting (GET & POST) dikunci khusus Admin
	sett := rg.Group("/setting", authMiddleware)
	{
		sett.GET("", handler.GetSetting)
		sett.POST("", handler.CreatedSettings)
		sett.POST("/update", handler.UpdateSetting)
	}
}
