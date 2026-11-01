package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"tjus-be/internal/config"
	"tjus-be/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(rdb *redis.Client, env *config.Environment) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: missing token")
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(env.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: invalid or expired token")
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: invalid claims")
			ctx.Abort()
			return
		}

		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: invalid user id claim")
			ctx.Abort()
			return
		}

		if _, err := uuid.Parse(userIdStr); err != nil {
			util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: invalid user id format")
			ctx.Abort()
			return
		}

		// Cek apakah token masuk daftar hitam (sudah Logout) di RAM Redis (0 query ke PostgreSQL!)
		if util.ExistsCache(ctx.Request.Context(), rdb, util.TokenBlacklistKey(tokenString)) {
			util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: session expired or logged out")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Set("email", claims["email"])

		// Set ke dalam Context GLOBAL agar bisa dipanggil dimana saja
		newCtx := context.WithValue(ctx.Request.Context(), "user_id", claims["user_id"])
		ctx.Request = ctx.Request.WithContext(newCtx)

		ctx.Next()
	}
}
