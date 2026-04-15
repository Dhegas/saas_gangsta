package errors

import (
	"errors"
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func Write(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.Status, appErr.Message, ErrorBody{
			Code:    appErr.Code,
			Details: appErr.Details,
		})
		return
	}

	response.Error(c, http.StatusInternalServerError, "Internal server error", ErrorBody{
		Code: "INTERNAL_ERROR",
	})
}

func Abort(c *gin.Context, err error) {
	Write(c, err)
	c.Abort()
}
