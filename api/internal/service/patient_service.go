package service

import (
	"github.com/SaidiBTW/appointment_booking_system_go/internal/domain"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository"
)

type PatientService interface {
	GetPatients() ([]*domain.Patient, error)
}

type patientService struct {
	patientRepository repository.PatientRepository
}

func NewPatientService(patientRepository repository.PatientRepository) PatientService {
	return &patientService{
		patientRepository: patientRepository,
	}
}

func (s *patientService) GetPatients() ([]*domain.Patient, error) {
	return s.patientRepository.GetPatients()
}
