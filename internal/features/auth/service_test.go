package auth

import (
	"context"
	"errors"
	"testing"
	"tjus-be/internal/config"
	"tjus-be/internal/pkg/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockAuthRepository adalah pemeran pengganti (Mock) untuk AuthRepository
type mockAuthRepository struct {
	mock.Mock
}

func (m *mockAuthRepository) FindByEmail(ctx context.Context, email string) (*AdminUser, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AdminUser), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte("rahasia123"), bcrypt.MinCost)
	require.NoError(t, err)

	dummyUser := &AdminUser{
		BaseEntity:   entity.BaseEntity{Id: uuid.New()},
		Email:        "admin@tjus.com",
		PasswordHash: string(hashedBytes),
		Name:         "Admin TJUS",
	}

	testEnv := &config.Environment{
		JwtSecret: "test-secret-key-minimal-32-karakter",
	}

	tests := []struct {
		name        string
		req         *LoginRequest
		setupMock   func(m *mockAuthRepository)
		wantErr     bool
		expectedErr string
	}{
		{
			name: "sukses login email dan password benar menghasilkan JWT token",
			req: &LoginRequest{
				Email:    "admin@tjus.com",
				Password: "rahasia123",
			},
			setupMock: func(m *mockAuthRepository) {
				m.On("FindByEmail", mock.Anything, "admin@tjus.com").Return(dummyUser, nil)
			},
			wantErr: false,
		},
		{
			name: "gagal login karena password salah",
			req: &LoginRequest{
				Email:    "admin@tjus.com",
				Password: "password_ngawur",
			},
			setupMock: func(m *mockAuthRepository) {
				m.On("FindByEmail", mock.Anything, "admin@tjus.com").Return(dummyUser, nil)
			},
			wantErr:     true,
			expectedErr: "invalid email or password",
		},
		{
			name: "gagal login karena user tidak ditemukan di database",
			req: &LoginRequest{
				Email:    "hantu@tjus.com",
				Password: "rahasia123",
			},
			setupMock: func(m *mockAuthRepository) {
				m.On("FindByEmail", mock.Anything, "hantu@tjus.com").Return(nil, errors.New("user not found"))
			},
			wantErr:     true,
			expectedErr: "invalid email or password",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockAuthRepository)
			tc.setupMock(mockRepo)

			authService := NewAuthService(mockRepo, testEnv, nil)

			resp, err := authService.Login(context.Background(), tc.req)

			if tc.wantErr {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.Token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

