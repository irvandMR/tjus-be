package asset

import (
	"net/http"
	"tjus-be/internal/pkg/httputil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	assetService *AssetService
}

func NewAssetHandler(assetService *AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

func (ah *AssetHandler) CreateAsset(c *gin.Context) {
	var request AssetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := ah.assetService.CreateAsset(c.Request.Context(), &request)
	if err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	httputil.ResponseSuccess(c, "Success Create Asset", resp)
}

func (ah *AssetHandler) GetAssetById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "invalid asset id format")
		return
	}

	resp, err := ah.assetService.GetAssetById(c.Request.Context(), id)
	if err != nil {
		httputil.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}

	httputil.ResponseSuccess(c, "Success Get Asset", resp)
}

func (ah *AssetHandler) GetAssetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		httputil.ResponseError(c, http.StatusBadRequest, "slug is required")
		return
	}

	resp, err := ah.assetService.GetAssetBySlug(c.Request.Context(), slug)
	if err != nil {
		httputil.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}

	httputil.ResponseSuccess(c, "Success Get Asset", resp)
}

func (ah *AssetHandler) GetAllAssets(c *gin.Context) {
	var filter AssetFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "invalid query parameters: "+err.Error())
		return
	}

	resp, pagination, err := ah.assetService.GetAllAssets(c.Request.Context(), filter)
	if err != nil {
		httputil.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.ResponsePagination(c, "Success Get Assets", resp, pagination)
}

func (ah *AssetHandler) GetAdminAssets(c *gin.Context) {
	var filter AssetFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "invalid query parameters: "+err.Error())
		return
	}

	resp, pagination, err := ah.assetService.GetAdminAssets(c.Request.Context(), filter)
	if err != nil {
		httputil.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.ResponsePagination(c, "Success Get Admin Assets", resp, pagination)
}

func (ah *AssetHandler) UploadMedia(c *gin.Context) {
	idParam := c.Param("id")
	assetId, err := uuid.Parse(idParam)
	if err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "invalid asset id format")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	// Mendukung pengiriman file dengan key "files" (multiple) atau "file" (single)
	files := form.File["files"]
	if len(files) == 0 {
		files = form.File["file"]
	}

	if len(files) == 0 {
		httputil.ResponseError(c, http.StatusBadRequest, "no files uploaded (use form key 'files' or 'file')")
		return
	}

	resp, err := ah.assetService.UploadMedia(c.Request.Context(), assetId, files)
	if err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	httputil.ResponseSuccess(c, "Success Upload Media", resp)
}

func (ah *AssetHandler) UpdateAsset(c *gin.Context) {
	var request UpdateAssetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := ah.assetService.UpdateAsset(c.Request.Context(), &request)
	if err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	httputil.ResponseSuccess(c, "Success Update Asset", resp)
}

func (ah *AssetHandler) DeleteMedia(c *gin.Context) {
	mediaIdParam := c.Param("media_id")
	mediaId, err := uuid.Parse(mediaIdParam)
	if err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "invalid media id format")
		return
	}

	if err := ah.assetService.DeleteMedia(c.Request.Context(), mediaId); err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	httputil.ResponseSuccess(c, "Success Delete Media", nil)
}

func (ah *AssetHandler) DeleteAsset(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, "invalid asset id format")
		return
	}

	if err := ah.assetService.DeleteAsset(c.Request.Context(), id); err != nil {
		httputil.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	httputil.ResponseSuccess(c, "Success Delete Asset", nil)
}
