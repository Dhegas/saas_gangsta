package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/gin-gonic/gin"
)

func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("panic recovered",
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("panic", fmt.Sprintf("%v", rec)),
				)
				apperrors.Abort(c, apperrors.New(
					"INTERNAL_ERROR",
					"An unexpected error occurred",
					http.StatusInternalServerError,
				))
			}
		}()

		c.Next()
	}
}
