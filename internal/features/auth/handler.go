package auth

import (
	"net/http"
	"tjus-be/internal/pkg/util"

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
		util.ResponseError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := a.authService.Login(c.Request.Context(), &req)
	if err != nil {
		util.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.ResponseSuccess(c, "Login success", resp)
}

func (a *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err := a.authService.Logout(c.Request.Context(), req.Token)

	if err != nil {
		util.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.Success(c, "Logout success")
}
