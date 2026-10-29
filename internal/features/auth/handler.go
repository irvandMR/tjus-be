package auth

import (
	"fmt"
	"net/http"
	"tjus-be/internal/pkg/httputil"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (a *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", a.Login)
	rg.POST("/logout", a.Logout)
}

func (a *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := a.authService.Login(c.Request.Context(), &req)
	fmt.Println("err : ", err)
	if err != nil {
		httputil.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.ResponseSuccess(c, "Login success", resp)
}

func (a *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err := a.authService.Logout(c.Request.Context(), req.Token)

	if err != nil {
		httputil.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.Success(c, "Logout success")
}
