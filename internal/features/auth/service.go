package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tjus-be/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository AuthRepository
	env            *config.Environment
}

func NewAuthService(authRepository AuthRepository, env *config.Environment) *AuthService {
	return &AuthService{authRepository: authRepository, env: env}
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

	err = a.authRepository.InsertToken(ctx, user.Id, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to insert token: %w", err)
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

	email, ok := claims["email"].(string)
	if !ok {
		return errors.New("email claim not found")
	}

	admin, err := a.authRepository.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if admin.Token == nil || *admin.Token != accessToken {
		return errors.New("token tidak valid atau sudah logout")
	}

	err = a.authRepository.DeleteToken(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil

}
