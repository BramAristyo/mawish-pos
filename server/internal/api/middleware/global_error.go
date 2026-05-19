package middleware

import (
	"errors"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/BramAristyo/mawish-pos/server/pkg/logger"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func ErrorHandler(log *logger.ZapLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		if duration > 200*time.Millisecond {
			log.Warn("Slow Request Detected",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"duration", duration.String(),
				"status", c.Writer.Status(),
			)
		}

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if errors.Is(err, io.EOF) {
				log.Warn("Request Body Empty", "path", c.Request.URL.Path, "error", err.Error())
				response.Error(c, http.StatusBadRequest, "request body is empty", err)
				return
			}

			// validation
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				log.Warn("Validation Failed", "path", c.Request.URL.Path, "fields", ve.Error())
				response.ValidationError(c, err)
				c.Errors = c.Errors[:0]
				return
			}

			// validation custom
			var customFieldErrs *usecase_errors.CustomFieldErrors
			if errors.As(err, &customFieldErrs) {
				log.Warn("Unique/Custom Field Violation", "path", c.Request.URL.Path)
				response.CustomValidationError(c, err)
				c.Errors = c.Errors[:0]
				return
			}

			var ue *usecase_errors.UseCaseError
			if errors.As(err, &ue) {
				if ue.Code >= 400 && ue.Code < 500 {
					log.Warn("Business Rule Violation",
						"path", c.Request.URL.Path,
						"code", ue.Code,
						"message", ue.Message,
					)
				} else {
					log.Error("Business Logic Error",
						"path", c.Request.URL.Path,
						"code", ue.Code,
						"error", err.Error(),
					)
				}
				response.Error(c, ue.Code, ue.Message, err)
				c.Errors = c.Errors[:0]
				return
			}

			// resource not found
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Warn("Resource Not Found", "path", c.Request.URL.Path, "error", err.Error())
				response.Error(c, http.StatusNotFound, "resource not found", err)
				return
			}

			// error duplicate
			if usecase_errors.IsUniqueViolation(err) {
				log.Warn("Duplicate Key Violation", "path", c.Request.URL.Path, "error", err.Error())
				response.Error(c, http.StatusConflict, "data already exists", err)
				return
			}

			log.Error("Critical System Error",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"error", err.Error(),
				"stack", string(debug.Stack()),
			)
			response.Error(c, http.StatusInternalServerError, "an unexpected error occurred", err)
		}
	}
}
