package repository_test

import (
	"testing"
	"time"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestAppointmentSlot_String(t *testing.T) {
	startTime, _ := time.Parse("15:04:05", "09:00:00")
	endTime, _ := time.Parse("15:04:05", "09:30:00")

	slot := repository.AppointmentSlot{
		StartTime: shared.TimeOnly(startTime),
		EndTime:   shared.TimeOnly(endTime),
	}

	assert.Equal(t, "09:00:00 - 09:30:00", slot.String())
}
