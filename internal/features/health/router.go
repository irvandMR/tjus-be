package health

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetHealthRouter(rg *gin.RouterGroup, db *pgxpool.Pool) {
	healthHandler := NewHealthHandler(db)
	healthHandler.RegisterRoutes(rg.Group("/health"))
}
