package service

import (
	"github.com/SaidiBTW/appointment_booking_system_go/internal/domain"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository"
)

type DoctorService interface {
	GetDoctors() ([]*domain.Doctor, error)
}

type doctorService struct {
	doctorRepository repository.DoctorRepository
}

func NewDoctorService(doctorRepository repository.DoctorRepository) DoctorService {
	return &doctorService{
		doctorRepository: doctorRepository,
	}
}

func (s *doctorService) GetDoctors() ([]*domain.Doctor, error) {
	return s.doctorRepository.GetDoctors()
}
