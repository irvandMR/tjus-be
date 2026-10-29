package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*AdminUser, error)
	FindTokenByUserId(ctx context.Context, id uuid.UUID) (*string, error)
	InsertToken(ctx context.Context, id uuid.UUID, token string) error
	DeleteToken(ctx context.Context, token string) error
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

func (a *authRepository) FindTokenByUserId(ctx context.Context, id uuid.UUID) (*string, error) {
	var token *string
	err := a.db.QueryRow(ctx, "SELECT token FROM tbl_admin_users WHERE id = $1", id).Scan(&token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return token, nil
}

func (a *authRepository) InsertToken(ctx context.Context, id uuid.UUID, token string) error {
	_, err := a.db.Exec(ctx, "UPDATE tbl_admin_users SET token = $1 WHERE id = $2", token, id)
	return err
}
func (a *authRepository) DeleteToken(ctx context.Context, token string) error {
	_, err := a.db.Exec(ctx, "UPDATE tbl_admin_users SET token = null WHERE token = $1", token)
	return err
}

func (a *authRepository) FindByEmail(ctx context.Context, email string) (*AdminUser, error) {
	var user AdminUser
	err := a.db.QueryRow(ctx, "SELECT id, email, password_hash, name, created_at, updated_at, token FROM tbl_admin_users WHERE email = $1", email).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt, &user.Token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
