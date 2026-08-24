package service

import (
	"time"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository"
)

type AvailabilityService interface {
	GetAvailableSlots(doctorID string, date time.Time) ([]repository.AppointmentSlot, error)
}

type availabilityService struct {
	availabilityRepository repository.AvailabilityRepository
}

func NewAvailabilityService(availabilityRepository repository.AvailabilityRepository) AvailabilityService {
	return &availabilityService{
		availabilityRepository: availabilityRepository,
	}
}

func (s *availabilityService) GetAvailableSlots(doctorID string, date time.Time) ([]repository.AppointmentSlot, error) {
	return s.availabilityRepository.GetAvailableSlots(doctorID, date)
}
