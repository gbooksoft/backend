package middleware

import (
	"errors"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

func ErrorHandler(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		var appErr gerr.AppError
		if !errors.As(err.Err, &appErr) {
			logger.Error().Err(err.Err).Msg("Unhandled error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		logger.Info().Err(appErr).Msg("Validation or business error")

		var validationErr *gerr.ValidationError
		if errors.As(appErr, &validationErr) {
			c.JSON(appErr.StatusCode(), gin.H{
				"error":          appErr.UserMessage(),
				"invalid_fields": validationErr.InvalidFields(),
			})
			return
		}

		c.JSON(appErr.StatusCode(), gin.H{
			"error":      appErr.UserMessage(),
			"error_code": appErr.PublicCode(),
		})
	}
}
