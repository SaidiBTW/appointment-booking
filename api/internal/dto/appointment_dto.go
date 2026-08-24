package dto

import "time"

type CreateAppointmentRequest struct {
	PatientID string    `json:"patient_id"`
	DoctorID  string    `json:"doctor_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time" validate:"gtfield=StartTime"`
}

type CreateAppointmentResponse struct {
	ID        string `json:"id"`
	PatientID string `json:"patient_id"`
	DoctorID  string `json:"doctor_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}

type RescheduleAppointmentRequest struct {
	NewStartTime time.Time `form:"new_start_time" time_format:"2006-01-02T15:04:05Z07:00" json:"new_start_time"`
	NewEndTime   time.Time `form:"new_end_time" time_format:"2006-01-02T15:04:05Z07:00" json:"new_end_time" validate:"gtfield=NewStartTime"`
}

type RescheduleAppointmentResponse struct {
	ID        string `json:"id"`
	PatientID string `json:"patient_id"`
	DoctorID  string `json:"doctor_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}
