package app

import (
	"tjus-be/internal/features/asset"
	"tjus-be/internal/features/auth"
	"tjus-be/internal/features/category"
	"tjus-be/internal/features/health"
	"tjus-be/internal/features/setting"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetRouter(s *Server) {
	// Serve static files for uploaded media
	s.Engine.Static("/uploads", "./uploads")

	version := s.Engine.Group("/api/v1")

	//Health
	health.SetHealthRouter(version, s.DB)

	//auth
	auth.SetAuthRouter(version, s.DB, s.Env)

	// Category
	category.SetCategoryRouter(version, s.DB, s.Env)

	// Setting
	setting.SetSettingRouting(version, s.DB, s.Env)

	// Asset
	asset.SetAssetRouter(version, s.DB, s.Env)

	authRepo := auth.NewAuthRepository(s.DB)
	authMiddleware := middleware.AuthMiddleware(authRepo, s.Env)
	version.GET("/test", authMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Selamat! Kamu lolos pemeriksaan satpam!"})
	})

	// Swagger Documentation & OpenAPI Spec
	s.Engine.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	s.Engine.StaticFile("/swagger", "./docs/index.html")
}
