package domain

import "time"

type Appointment struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patient_id"`
	DoctorID  string    `json:"doctor_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}
