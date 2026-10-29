package setting

import (
	"tjus-be/internal/config"
	"tjus-be/internal/features/auth"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetSettingRouting(rg *gin.RouterGroup, db *pgxpool.Pool, env *config.Environment) {
	repository := NewSettingsRepository(db)
	service := NewSettingsService(repository)
	handler := NewSettingsHandler(service)

	authRepo := auth.NewAuthRepository(db)
	authMiddleware := middleware.AuthMiddleware(authRepo, env)

	// Semua endpoint setting (GET & POST) dikunci khusus Admin
	sett := rg.Group("/setting", authMiddleware)
	{
		sett.GET("", handler.GetSetting)
		sett.POST("", handler.CreatedSettings)
		sett.POST("/update", handler.UpdateSetting)
	}
}
