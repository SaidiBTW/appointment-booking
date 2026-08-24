package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/domain"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/dto"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/service"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/shared"
	"github.com/gin-gonic/gin"
)

type MockAppointmentService interface {
	CreateAppointment(appointmentDto dto.CreateAppointmentRequest) (*domain.Appointment, error)
	CancelAppointment(appointmentID string, patientID string, reason string) error
	RescheduleAppointment(appointmentID string, newStartTime time.Time, newEndTime time.Time) (*domain.Appointment, error)
	GetAppointmentsByPatientID(patientID string) ([]*domain.Appointment, error)
}

type MockAvailabilityService interface {
	GetAvailableSlots(doctorID string, date time.Time) ([]string, error)
}

func NewMockAppointmentService() service.AppointmentService {
	return &mockAppointmentService{}
}

func NewMockAvailabilityService() service.AvailabilityService {
	return &mockAvailabilityService{}
}

type mockAvailabilityService struct{}

func (m *mockAvailabilityService) GetAvailableSlots(doctorID string, date time.Time) ([]repository.AppointmentSlot, error) {
	// Mock implementation for getting available slots
	return []repository.AppointmentSlot{
		{
			StartTime: shared.TimeOnly(date.Add(9 * time.Hour)),
			EndTime:   shared.TimeOnly(date.Add(9*time.Hour + 30*time.Minute)),
		},
		{
			StartTime: shared.TimeOnly(date.Add(10 * time.Hour)),
			EndTime:   shared.TimeOnly(date.Add(10*time.Hour + 30*time.Minute)),
		},
	}, nil
}

type mockAppointmentService struct{}

func (m *mockAppointmentService) CreateAppointment(appointmentDto dto.CreateAppointmentRequest) (*domain.Appointment, error) {
	// Mock implementation for creating an appointment
	return &domain.Appointment{
		ID:        "mock-id",
		DoctorID:  appointmentDto.DoctorID,
		PatientID: appointmentDto.PatientID,
		StartTime: appointmentDto.StartTime,
		EndTime:   appointmentDto.EndTime,
		Status:    "Scheduled",
	}, nil
}

func (m *mockAppointmentService) CancelAppointment(appointmentID string, patientID string, reason string) error {
	// Mock implementation for canceling an appointment
	return nil
}

func (m *mockAppointmentService) RescheduleAppointment(appointmentID string, newStartTime time.Time, newEndTime time.Time) (*domain.Appointment, error) {
	// Mock implementation for rescheduling an appointment
	return &domain.Appointment{
		ID:        appointmentID,
		DoctorID:  uuid.New().String(),
		PatientID: uuid.New().String(),
		StartTime: newStartTime,
		EndTime:   newEndTime,
		Status:    "Scheduled",
	}, nil
}

func (m *mockAppointmentService) GetAppointmentsByPatientID(patientID string) ([]*domain.Appointment, error) {
	// Mock implementation for getting appointments by patient ID
	return []*domain.Appointment{
		{
			ID:        uuid.New().String(),
			DoctorID:  uuid.New().String(),
			PatientID: patientID,
			StartTime: time.Now().Add(24 * time.Hour),
			EndTime:   time.Now().Add(25 * time.Hour),
			Status:    "Scheduled",
		},
	}, nil
}

func TestAppointmentHandler_HandleCreateAppointment(t *testing.T) {

	gin.SetMode(gin.TestMode)

	appointmentService := NewMockAppointmentService()
	appointmentHandler := NewAppointmentHandler(appointmentService)

	router := gin.Default()
	router.POST("/appointments", appointmentHandler.CreateAppointment)

	t.Run("fails within 1 hour creation", func(t *testing.T) {
		// Test logic for failing within 1 hour creation
		doctor_uuid := uuid.New().String()
		patient_uuid := uuid.New().String()
		startTime := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
		endTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		payload := `{
		"doctor_id": "` + doctor_uuid + `",
		"patient_id": "` + patient_uuid + `",
		"start_time": "` + startTime + `",
		"end_time": "` + endTime + `"
	}`

		req, _ := http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("only 30 minute blocks", func(t *testing.T) {
		doctor_uuid := uuid.New().String()
		patient_uuid := uuid.New().String()
		startTime := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(3 * time.Hour).Format(time.RFC3339)
		payload := `{
		"doctor_id": "` + doctor_uuid + `",
		"patient_id": "` + patient_uuid + `",
		"start_time": "` + startTime + `",
		"end_time": "` + endTime + `"
	}`

		req, _ := http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("successfully create appointment", func(t *testing.T) {
		doctor_uuid := uuid.New().String()
		patient_uuid := uuid.New().String()
		startTime := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(2*time.Hour + 30*time.Minute).Format(time.RFC3339)
		payload := `{
		"doctor_id": "` + doctor_uuid + `",
		"patient_id": "` + patient_uuid + `",
		"start_time": "` + startTime + `",
		"end_time": "` + endTime + `"
	}`

		req, _ := http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})
}

func TestAppointmentHandler_HandleCancelAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	appointmentService := NewMockAppointmentService()
	appointmentHandler := NewAppointmentHandler(appointmentService)

	router := gin.Default()
	router.DELETE("/appointments/:id", appointmentHandler.CancelAppointment)

	t.Run("successfully cancel appointment", func(t *testing.T) {
		appointmentID := "mock-id"
		patientID := uuid.New().String()
		reason := "Patient is unavailable"

		req, _ := http.NewRequest(http.MethodDelete, "/appointments/"+appointmentID+"?patient_id="+patientID+"&reason="+reason, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestAppointmentHandler_HandleRescheduleAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	appointmentService := NewMockAppointmentService()
	appointmentHandler := NewAppointmentHandler(appointmentService)

	router := gin.Default()
	router.PATCH("/appointments/:id/reschedule", appointmentHandler.RescheduleAppointment)

	t.Run("successfully reschedule appointment", func(t *testing.T) {
		appointmentID := uuid.New().String()
		newStartTime := time.Now().Add(3 * time.Hour).Format(time.RFC3339)
		newEndTime := time.Now().Add(3*time.Hour + 30*time.Minute).Format(time.RFC3339)

		q := url.Values{}
		q.Add("new_start_time", newStartTime)
		q.Add("new_end_time", newEndTime)

		urlPath := "/appointments/" + appointmentID + "/reschedule?" + q.Encode()

		req, err := http.NewRequest(http.MethodPatch, urlPath, nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d: %s: %v ", http.StatusOK, w.Code, w.Body.String(), newStartTime)
		}
	})

	// 1 hour reschedule restriction test
	t.Run("fails within 1 hour reschedule", func(t *testing.T) {
		appointmentID := uuid.New().String()
		newStartTime := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
		newEndTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

		q := url.Values{}
		q.Add("new_start_time", newStartTime)
		q.Add("new_end_time", newEndTime)

		urlPath := "/appointments/" + appointmentID + "/reschedule?" + q.Encode()

		req, err := http.NewRequest(http.MethodPatch, urlPath, nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestAppointmentHandler_HandleGetAppointmentsByPatientID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	appointmentService := NewMockAppointmentService()
	appointmentHandler := NewAppointmentHandler(appointmentService)

	router := gin.Default()
	router.GET("/appointments/patient/:id", appointmentHandler.GetAppointmentsByPatientID)

	t.Run("successfully get appointments by patient ID", func(t *testing.T) {
		patientID := uuid.New().String()

		req, _ := http.NewRequest(http.MethodGet, "/appointments/patient/"+patientID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

}
func TestAvailabilityHandler_HandleGetAvailability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	availabilityService := NewMockAvailabilityService()
	availabilityHandler := NewAvailabilityHandler(availabilityService)

	router := gin.Default()
	router.GET("/availability/:id", availabilityHandler.GetAvailability)

	t.Run("successfully get available slots", func(t *testing.T) {
		doctorID := uuid.New().String()
		date := time.Now().Add(24 * time.Hour).Format("2006-01-02")

		q := url.Values{}
		q.Add("date", date)

		urlPath := "/availability/" + doctorID + "?" + q.Encode()

		req, _ := http.NewRequest(http.MethodGet, urlPath, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}
