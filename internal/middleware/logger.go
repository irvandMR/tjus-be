package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware mencatat aktivitas HTTP secara terstruktur (Method, Path, Status, Latency, IP, dan User-Agent).
func LoggerMiddleware(isProduction bool) gin.HandlerFunc {
	var handler slog.Handler
	if isProduction {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(handler)

	return func(c *gin.Context) {
		// 1. Pintu Masuk (Sebelum c.Next): Nyalakan stopwatch & catat URL
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		// 2. Teruskan request ke Handler -> Service -> Redis/PostgreSQL
		c.Next()

		// 3. Pintu Keluar (Sesudah c.Next): Hitung durasi & kumpulkan informasi response
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		attrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.String("latency", latency.String()),
			slog.String("ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		// Jika request berasal dari Admin yang sedang login, ikut catat email-nya (Audit Trail)
		if email, exists := c.Get("email"); exists {
			if emailStr, ok := email.(string); ok && emailStr != "" {
				attrs = append(attrs, slog.String("admin_email", emailStr))
			}
		}

		// 4. Tentukan level keparahan log berdasarkan kode HTTP Status
		switch {
		case statusCode >= 500:
			logger.Error("HTTP Request", attrs...)
		case statusCode >= 400:
			logger.Warn("HTTP Request", attrs...)
		default:
			logger.Info("HTTP Request", attrs...)
		}
	}
}

