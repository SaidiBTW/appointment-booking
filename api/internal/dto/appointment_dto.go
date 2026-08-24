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
	NewStartTime string `json:"new_start_time"`
	NewEndTime   string `json:"new_end_time" validate:"gtefield=NewStartTime"`
}

type RescheduleAppointmentResponse struct {
	ID        string `json:"id"`
	PatientID string `json:"patient_id"`
	DoctorID  string `json:"doctor_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}
