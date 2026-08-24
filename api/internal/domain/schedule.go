package domain

import "github.com/SaidiBTW/appointment_booking_system_go/internal/shared"

type Schedule struct {
	ID        string          `json:"id"`
	DoctorID  string          `json:"doctor_id"`
	DayOfWeek int             `json:"day_of_week"`
	StartTime shared.TimeOnly `json:"start_time"`
	EndTime   shared.TimeOnly `json:"end_time"`
}
