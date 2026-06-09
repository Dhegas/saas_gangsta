package errors

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/gin-gonic/gin"
)

// ErrorBody adalah struktur field "error" dalam response standar.
type ErrorBody struct {
	Code string `json:"code"`
}

// ValidationErrorBody digunakan khusus untuk error validasi dengan field-level errors.
type ValidationErrorBody struct {
	Code string `json:"code"`
}

// Write menulis error response standar ke client.
// Error internal di-log ke server; client hanya menerima pesan aman.
func Write(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.Status, appErr.Message, ErrorBody{
			Code: appErr.Code,
		})
		return
	}

	// Unexpected error: log detail internal, kembalikan respons generik
	slog.Error("unexpected error",
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.String("error", err.Error()),
	)

	response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", ErrorBody{
		Code: "INTERNAL_ERROR",
	})
}

// WriteValidation menulis error validasi dengan format field-level errors.
func WriteValidation(c *gin.Context, message string, fieldErrors map[string][]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"success": false,
		"message": message,
		"error": ValidationErrorBody{
			Code: "VALIDATION_ERROR",
		},
		"errors": fieldErrors,
	})
}

// Abort memanggil Write lalu menghentikan chain middleware.
func Abort(c *gin.Context, err error) {
	Write(c, err)
	c.Abort()
}
