package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*AdminUser, error)
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

func (a *authRepository) FindByEmail(ctx context.Context, email string) (*AdminUser, error) {
	var user AdminUser
	err := a.db.QueryRow(
		ctx,
		"SELECT id, email, password_hash, name, created_at, updated_at, token FROM tbl_admin_users WHERE email = $1",
		email,
	).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt, &user.Token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
