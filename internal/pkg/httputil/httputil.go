package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

type PaginatedResponse struct {
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func ResponseSuccess(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, BaseResponse{
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	})
}

func ResponsePagination(c *gin.Context, message string, data any, pagination Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func Success(c *gin.Context, message string) {
	c.JSON(http.StatusOK, BaseResponse{
		StatusCode: http.StatusOK,
		Message:    message,
	})
}

func ResponseError(c *gin.Context, code int, message string) {
	c.JSON(code, BaseResponse{
		StatusCode: code,
		Message:    message,
		Data:       nil,
	})
}

