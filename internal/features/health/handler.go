package health

import (
	"net/http"
	"tjus-be/internal/pkg/httputil"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Method untuk pasang rute health ke router Gin
func (h *HealthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetHealthHandler)
	rg.GET("/db", h.GetHealthHandlerDB)
}

func (h *HealthHandler) GetHealthHandler(c *gin.Context) {
	httputil.ResponseSuccess(c, "Server is running", nil)
}

func (h *HealthHandler) GetHealthHandlerDB(c *gin.Context) {
	err := h.db.Ping(c.Request.Context())
	if err != nil {
		httputil.ResponseError(c, http.StatusInternalServerError, "Database is unreachable")
		return
	}
	httputil.ResponseSuccess(c, "Database is healthy", nil)
}
