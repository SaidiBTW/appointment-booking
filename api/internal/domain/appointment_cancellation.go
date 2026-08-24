package domain

type AppointmentCancellation struct {
	ID            string `json:"id"`
	AppointmentID string `json:"appointment_id"`
	PatientID     string `json:"patient_id"`
	Reason        string `json:"reason"`
}
