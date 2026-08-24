package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/SaidiBTW/appointment_booking_system_go/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	return r
}

func TestErrorHandler_Success(t *testing.T) {
	router := setupTestRouter()
	router.GET("/success", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/success", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

func TestErrorHandler_APIError(t *testing.T) {
	router := setupTestRouter()
	router.GET("/bad-request", func(c *gin.Context) {
		_ = c.Error(errors.BadRequestError("Invalid query parameter", fmt.Errorf("missing id")))
	})

	req, _ := http.NewRequest(http.MethodGet, "/bad-request", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{
		"success": false,
		"code": "bad_request",
		"message": "Invalid query parameter"
	}`, w.Body.String())
}

func TestErrorHandler_InternalServerError(t *testing.T) {
	router := setupTestRouter()
	router.GET("/generic-error", func(c *gin.Context) {
		_ = c.Error(fmt.Errorf("something unexpected broke"))
	})

	req, _ := http.NewRequest(http.MethodGet, "/generic-error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{
		"success": false,
		"code": "internal_server_error",
		"message": "An internal server error occurred"
	}`, w.Body.String())
}
