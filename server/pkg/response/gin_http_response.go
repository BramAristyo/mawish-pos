package response

import (
	"net/http"

	"github.com/BramAristyo/mawish-pos/server/internal/api/validation"
	"github.com/gin-gonic/gin"
)

type BaseHTTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func OKPaginate[T any, M any](c *gin.Context, data T, meta M) {
	c.JSON(http.StatusOK, BaseHTTPResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func OK(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, BaseHTTPResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Created(c *gin.Context, data any, message string) {
	c.JSON(http.StatusCreated, BaseHTTPResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Error(c *gin.Context, code int, message string, err error) {
	c.JSON(code, BaseHTTPResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
}

func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, BaseHTTPResponse{
		Success: false,
		Error:   validation.GetValidationErrors(err),
	})
}

func CustomValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "Validation Failed",
		"error":   err,
	})

	c.Errors = c.Errors[:0]
}
