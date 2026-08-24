package repository

import (
	"context"
	"database/sql"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/domain"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository/gen_queries"
)

type PatientRepository interface {
	GetPatients() ([]*domain.Patient, error)
}

type patientRepository struct {
	db *sql.DB
}

func NewPostgresPatientRepository(db *sql.DB) PatientRepository {
	return &patientRepository{
		db: db,
	}
}

func (r *patientRepository) GetPatients() ([]*domain.Patient, error) {
	ctx := context.Background()
	queries := gen_queries.New(r.db)
	query, err := queries.GetPatients(ctx)
	if err != nil {
		return nil, err
	}

	var patients []*domain.Patient

	for _, p := range query {
		patient := &domain.Patient{
			ID:   p.ID.String(),
			Name: p.Name,
		}
		patients = append(patients, patient)
	}

	return patients, nil
}
