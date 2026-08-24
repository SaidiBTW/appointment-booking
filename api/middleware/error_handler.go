package middleware

import (
	"net/http"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.StatusCode, gin.H{
				"success": false,
				"code":    apiErr.Code,
				"message": apiErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "internal_server_error",
			"message": "An internal server error occurred",
		})

	}
}
