package category

import (
	"tjus-be/internal/config"
	"tjus-be/internal/features/auth"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetCategoryRouter(rg *gin.RouterGroup, db *pgxpool.Pool, env *config.Environment) {
	repository := NewCategoryRepository(db)
	service := NewCategoryService(repository)
	handler := NewCategoryHandler(service)

	authRepo := auth.NewAuthRepository(db)
	authMiddleware := middleware.AuthMiddleware(authRepo, env)

	ctg := rg.Group("/category")
	{
		// Categories
		ctg.GET("", handler.GetAllCategories)
		ctg.POST("", authMiddleware, handler.CreatedCatagory)
		ctg.POST("/update", authMiddleware, handler.EditCatagory)
		ctg.DELETE("/:id", authMiddleware, handler.DeleteCategory)

		// Category Attributes
		ctg.POST("/attribute", authMiddleware, handler.CreateAttribute)
		ctg.POST("/attribute/update", authMiddleware, handler.EditAttribute)
		ctg.GET("/:id/attributes", handler.GetAttributesByCategoryId)
		ctg.DELETE("/attribute/:attr_id", authMiddleware, handler.DeleteAttribute)
	}

}
