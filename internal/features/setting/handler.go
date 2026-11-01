package setting

import (
	"net/http"
	"tjus-be/internal/pkg/util"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsService *SettingsService
}

func NewSettingsHandler(settingsService *SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: settingsService}
}

func (sh *SettingsHandler) CreatedSettings(c *gin.Context) {
	var request SettingsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	err := sh.settingsService.CreateSettings(c.Request.Context(), &request)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	util.Success(c, "Success Create Setting")
}

func (sh *SettingsHandler) GetSetting(c *gin.Context) {
	resp, err := sh.settingsService.GetSetting(c.Request.Context())
	if err != nil {
		util.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.ResponseSuccess(c, "Success Get Setting", resp)
}

func (sh *SettingsHandler) UpdateSetting(c *gin.Context) {
	var request UpdateSettingsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := sh.settingsService.EditSettings(c.Request.Context(), &request)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Success Update Setting", resp)
}
