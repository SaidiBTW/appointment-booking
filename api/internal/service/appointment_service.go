package service

import (
	"time"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/domain"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/dto"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository"
)

type AppointmentService interface {
	CreateAppointment(appointmentDto dto.CreateAppointmentRequest) (*domain.Appointment, error)
	CancelAppointment(appointmentID string, patientID string, reason string) error
	RescheduleAppointment(appointmentID string, newStartTime time.Time, newEndTime time.Time) (*domain.Appointment, error)
	GetAppointmentsByPatientID(patientID string) ([]*domain.Appointment, error)
}

type appointmentService struct {
	appointmentRepository repository.AppointmentRepository
}

func NewAppointmentService(appointmentRepository repository.AppointmentRepository) AppointmentService {
	return &appointmentService{
		appointmentRepository: appointmentRepository,
	}
}

func (s *appointmentService) CreateAppointment(appointmentDto dto.CreateAppointmentRequest) (*domain.Appointment, error) {

	return s.appointmentRepository.CreateAppointment(appointmentDto)
}

func (s *appointmentService) CancelAppointment(appointmentID string, patientID string, reason string) error {
	return s.appointmentRepository.CancelAppointment(appointmentID, patientID, reason)
}

func (s *appointmentService) RescheduleAppointment(appointmentID string, newStartTime time.Time, newEndTime time.Time) (*domain.Appointment, error) {
	return s.appointmentRepository.RescheduleAppointment(appointmentID, newStartTime, newEndTime)
}

func (s *appointmentService) GetAppointmentsByPatientID(patientID string) ([]*domain.Appointment, error) {
	return s.appointmentRepository.GetAppointmentsByPatientID(patientID)
}
