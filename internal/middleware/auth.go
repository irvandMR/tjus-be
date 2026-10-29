package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"tjus-be/internal/config"
	"tjus-be/internal/features/auth"
	"tjus-be/internal/pkg/httputil"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(authRepo auth.AuthRepository, env *config.Environment) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			httputil.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: missing token")
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
			httputil.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: invalid or expired token")
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			httputil.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: invalid claims")
			ctx.Abort()
			return
		}

		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			httputil.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: invalid user id claim")
			ctx.Abort()
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			httputil.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: invalid user id format")
			ctx.Abort()
			return
		}

		// Pengecekan Stateful Token lewat AuthRepository (query ada di repo)
		activeToken, err := authRepo.FindTokenByUserId(ctx.Request.Context(), userId)
		if err != nil || activeToken == nil || *activeToken != tokenString {
			httputil.ResponseError(ctx, http.StatusUnauthorized, "unauthorized: session expired or logged out")
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
