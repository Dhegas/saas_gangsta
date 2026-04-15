package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
	c.JSON(status, Envelope{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, status int, message string, errBody interface{}) {
	c.JSON(status, Envelope{
		Success: false,
		Message: message,
		Error:   errBody,
	})
}
