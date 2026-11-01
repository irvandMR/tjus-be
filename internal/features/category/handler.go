package category

import (
	"net/http"
	"tjus-be/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService *CategoryService
}

func NewCategoryHandler(categoryService *CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (ch *CategoryHandler) CreatedCatagory(c *gin.Context) {
	var request CreatedCatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := ch.categoryService.CreateCategory(c.Request.Context(), &request)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Success Create Category", resp)
}

func (ch *CategoryHandler) EditCatagory(c *gin.Context) {
	var request UpdateCatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := ch.categoryService.EditCategory(c.Request.Context(), &request)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Success Update Category", resp)
}

func (ch *CategoryHandler) GetAllCategories(c *gin.Context) {
	resp, err := ch.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		util.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.ResponseSuccess(c, "Success Get Categories", resp)
}

func (ch *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid category id format")
		return
	}

	err = ch.categoryService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Success Delete Category", nil)
}

func (ch *CategoryHandler) CreateAttribute(c *gin.Context) {
	var request CreatedCatAttrRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := ch.categoryService.CreateAttribute(c.Request.Context(), &request)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Success Create Category Attributes", resp)
}

func (ch *CategoryHandler) GetAttributesByCategoryId(c *gin.Context) {
	idStr := c.Param("id")
	categoryId, err := uuid.Parse(idStr)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid category id format")
		return
	}

	resp, err := ch.categoryService.GetAttributesByCategoryId(c.Request.Context(), categoryId)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Success Get Category Attributes", resp)
}

func (ch *CategoryHandler) EditAttribute(c *gin.Context) {
	var request UpdateCatAttrRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := ch.categoryService.EditAttribute(c.Request.Context(), &request)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Success Update Category Attribute", resp)
}

func (ch *CategoryHandler) DeleteAttribute(c *gin.Context) {
	idStr := c.Param("attr_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, "invalid attribute id format")
		return
	}

	err = ch.categoryService.DeleteAttribute(c.Request.Context(), id)
	if err != nil {
		util.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Success Delete Category Attribute", nil)
}
