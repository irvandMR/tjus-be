package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tjus-be/internal/config"
	"tjus-be/internal/pkg/entity"
	"tjus-be/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte("rahasia123"), bcrypt.MinCost)
	require.NoError(t, err)

	dummyUser := &AdminUser{
		BaseEntity:   entity.BaseEntity{Id: uuid.New()},
		Email:        "admin@tjus.com",
		PasswordHash: string(hashedBytes),
		Name:         "Admin TJUS",
	}

	testEnv := &config.Environment{JwtSecret: "test-secret-key-minimal-32-karakter"}

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(m *mockAuthRepository)
		expectedStatus int
		expectedMsg    string
		expectToken    bool
	}{
		{
			name:        "200 OK saat login berhasil",
			requestBody: `{"email": "admin@tjus.com", "password": "rahasia123"}`,
			setupMock: func(m *mockAuthRepository) {
				m.On("FindByEmail", mock.Anything, "admin@tjus.com").Return(dummyUser, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Login success",
			expectToken:    true,
		},
		{
			name:           "400 Bad Request saat payload JSON kosong",
			requestBody:    `{}`,
			setupMock:      func(m *mockAuthRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Invalid request payload",
			expectToken:    false,
		},
		{
			name:        "401 Unauthorized saat password salah",
			requestBody: `{"email": "admin@tjus.com", "password": "salah"}`,
			setupMock: func(m *mockAuthRepository) {
				m.On("FindByEmail", mock.Anything, "admin@tjus.com").Return(dummyUser, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "invalid email or password",
			expectToken:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockAuthRepository)
			tc.setupMock(mockRepo)

			service := NewAuthService(mockRepo, testEnv, nil)
			handler := NewAuthHandler(service)

			router := gin.New()
			handler.RegisterRoutes(router.Group("/api/v1/auth"))

			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(tc.requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatus, recorder.Code)

			var body util.BaseResponse
			err = json.Unmarshal(recorder.Body.Bytes(), &body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, body.StatusCode)
			assert.Equal(t, tc.expectedMsg, body.Message)

			if tc.expectToken {
				dataMap, ok := body.Data.(map[string]any)
				require.True(t, ok)
				assert.NotEmpty(t, dataMap["token"])
			} else {
				assert.Nil(t, body.Data)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

