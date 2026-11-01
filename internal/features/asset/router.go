package asset

import (
	"tjus-be/internal/config"
	"tjus-be/internal/features/category"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetAssetRouter(rg *gin.RouterGroup, db *pgxpool.Pool, env *config.Environment, rdb *redis.Client) {
	repository := NewAssetRepository(db)
	catRepository := category.NewCategoryRepository(db)
	service := NewAssetService(repository, catRepository, rdb)
	handler := NewAssetHandler(service)

	authMiddleware := middleware.AuthMiddleware(rdb, env)

	assetGroup := rg.Group("/assets")
	{
		// Public routes (Katalog & Detail untuk pengunjung - Tanpa data Owner)
		assetGroup.GET("", handler.GetAllAssets)
		assetGroup.GET("/slug/:slug", handler.GetAssetBySlug)

		// Protected routes (Admin Dashboard - Wajib bawa token & ada data Owner)
		assetGroup.GET("/admin", authMiddleware, handler.GetAdminAssets)
		assetGroup.GET("/:id", authMiddleware, handler.GetAssetById)
		assetGroup.POST("", authMiddleware, handler.CreateAsset)
		assetGroup.POST("/update", authMiddleware, handler.UpdateAsset)
		assetGroup.POST("/:id/media", authMiddleware, handler.UploadMedia)
		assetGroup.DELETE("/media/:media_id", authMiddleware, handler.DeleteMedia)
		assetGroup.DELETE("/:id", authMiddleware, handler.DeleteAsset)
	}
}
