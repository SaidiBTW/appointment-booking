package errors_test

import (
	"fmt"
	"testing"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestAPIError_ErrorFormatting(t *testing.T) {
	t.Run("Format with inner error", func(t *testing.T) {
		err := errors.New(400, "invalid_payload", "Payload is invalid", fmt.Errorf("field missing"))
		assert.Equal(t, "[invalid_payload] Payload is invalid: field missing", err.Error())
	})

	t.Run("Format without inner error", func(t *testing.T) {
		err := errors.CustomError(404, "not_found", "Resource not found")
		assert.Equal(t, "[not_found] Resource not found", err.Error())
	})
}

func TestAPIError_Constructors(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := errors.NotFoundError("Doctor not found")
		assert.Equal(t, 404, err.StatusCode)
		assert.Equal(t, "not_found", err.Code)
		assert.Equal(t, "Doctor not found", err.Message)
		assert.Nil(t, err.Err)
	})

	t.Run("BadRequestError", func(t *testing.T) {
		innerErr := fmt.Errorf("invalid uuid")
		err := errors.BadRequestError("Bad patient ID", innerErr)
		assert.Equal(t, 400, err.StatusCode)
		assert.Equal(t, "bad_request", err.Code)
		assert.Equal(t, "Bad patient ID", err.Message)
		assert.Equal(t, innerErr, err.Err)
	})

}
