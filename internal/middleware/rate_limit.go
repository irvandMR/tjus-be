package middleware

import (
	"fmt"
	"net/http"
	"time"
	"tjus-be/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitMiddleware membatasi jumlah request per IP dalam rentang waktu tertentu menggunakan Redis.
func RateLimitMiddleware(rdb *redis.Client, prefix string, maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s:%s", prefix, c.ClientIP())

		count, err := util.IncrementCounter(c.Request.Context(), rdb, key, window)
		if err != nil {
			// Fail-open jika koneksi Redis gangguan sesaat
			c.Next()
			return
		}

		if count > maxRequests {
			util.ResponseError(c, http.StatusTooManyRequests, "too many requests, please try again later")
			c.Abort()
			return
		}

		c.Next()
	}
}
