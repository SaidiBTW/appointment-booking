package repository

import (
	"context"
	"database/sql"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/domain"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository/gen_queries"
)

type DoctorRepository interface {
	GetDoctors() ([]*domain.Doctor, error)
}

type doctorRepository struct {
	db *sql.DB
}

func NewPostgresDoctorRepository(db *sql.DB) DoctorRepository {
	return &doctorRepository{
		db: db,
	}
}

func (r *doctorRepository) GetDoctors() ([]*domain.Doctor, error) {
	ctx := context.Background()
	queries := gen_queries.New(r.db)
	query, err := queries.GetDoctors(ctx)
	if err != nil {
		return nil, err
	}

	var doctors []*domain.Doctor

	for _, d := range query {
		doctor := &domain.Doctor{
			ID:   d.ID.String(),
			Name: d.Name,
		}
		doctors = append(doctors, doctor)
	}

	return doctors, nil
}
