package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tjus-be/internal/config"
	"tjus-be/internal/pkg/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository AuthRepository
	env            *config.Environment
	redisClient    *redis.Client
}

func NewAuthService(authRepository AuthRepository, env *config.Environment, redisClient *redis.Client) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		env:            env,
		redisClient:    redisClient,
	}
}

func (a *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := a.authRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
		"iss":     "tjus-be",
		"sub":     "access_token",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.env.JwtSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{
		Token: tokenString,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, accessToken string) error {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.env.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return errors.New("token tidak valid atau sudah kadaluarsa")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	if a.redisClient == nil {
		return errors.New("redis sedang tidak aktif, gagal memproses logout")
	}

	blacklistKey := util.TokenBlacklistKey(accessToken)
	if util.ExistsCache(ctx, a.redisClient, blacklistKey) {
		return errors.New("token tidak valid atau sudah logout")
	}

	expTime, err := claims.GetExpirationTime()
	if err != nil || expTime == nil {
		return errors.New("invalid token expiration claim")
	}

	remainingTTL := time.Until(expTime.Time)
	if remainingTTL <= 0 {
		return errors.New("token sudah kadaluarsa")
	}

	if err := util.SetCache(ctx, a.redisClient, blacklistKey, "revoked", remainingTTL); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}
