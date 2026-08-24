package repository

import (
	"context"
	"database/sql"
	"log"

	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/domain"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/dto"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository/gen_queries"
)

type AppointmentRepository interface {
	CreateAppointment(appointmentDto dto.CreateAppointmentRequest) (*domain.Appointment, error)
	GetAppointmentByID(id string) (*domain.Appointment, error)
	GetAppointmentsByDoctorID(doctorID string) ([]*domain.Appointment, error)
	GetAppointmentsByPatientID(patientID string) ([]*domain.Appointment, error)
	RescheduleAppointment(appointmentID string, newStartTime time.Time, newEndTime time.Time) (*domain.Appointment, error)
	CancelAppointment(appointmentID string, patientID string, reason string) error
}

type appointmentRepository struct {
	db *sql.DB
}

func NewPostgresAppointmentRepository(db *sql.DB) AppointmentRepository {

	return &appointmentRepository{
		db: db,
	}
}

func (r *appointmentRepository) CreateAppointment(appointmentDto dto.CreateAppointmentRequest) (*domain.Appointment, error) {
	ctx := context.Background()
	queries := gen_queries.New(r.db)
	// Within 1 hour rejection logic
	log.Printf("Time in DTO is %v", appointmentDto.StartTime)
	query, err := queries.CreateAppointment(ctx, gen_queries.CreateAppointmentParams{
		PatientID: uuid.NullUUID{UUID: uuid.MustParse(appointmentDto.PatientID), Valid: true},
		DoctorID:  uuid.NullUUID{UUID: uuid.MustParse(appointmentDto.DoctorID), Valid: true},
		StartTime: appointmentDto.StartTime,
		EndTime:   appointmentDto.EndTime,
		Status:    "scheduled",
	})

	if err != nil {
		return nil, err
	}
	var appointment domain.Appointment

	appointment.ID = query.ID.String()
	appointment.PatientID = query.PatientID.UUID.String()
	appointment.DoctorID = query.DoctorID.UUID.String()
	appointment.StartTime = query.StartTime
	appointment.EndTime = query.EndTime
	appointment.Status = query.Status
	return &appointment, nil
}

func (r *appointmentRepository) GetAppointmentByID(id string) (*domain.Appointment, error) {
	ctx := context.Background()
	queries := gen_queries.New(r.db)
	query, err := queries.GetAppointmentByID(ctx,
		uuid.MustParse(id),
	)
	if err != nil {
		return nil, err
	}

	var appointment domain.Appointment
	appointment.ID = query.ID.String()
	appointment.PatientID = query.PatientID.UUID.String()
	appointment.DoctorID = query.DoctorID.UUID.String()
	appointment.StartTime = query.StartTime
	appointment.EndTime = query.EndTime
	appointment.Status = query.Status

	return &appointment, nil

}

func (r *appointmentRepository) RescheduleAppointment(appointmentID string, newStartTime time.Time, newEndTime time.Time) (*domain.Appointment, error) {
	ctx := context.Background()
	queries := gen_queries.New(r.db)

	appointmentsInRange, err := r.GetAppointmentsByDoctorIDAndStartAndEndTime(appointmentID, newStartTime, newEndTime)
	if err != nil {
		return nil, err
	}
	if len(appointmentsInRange) > 0 {
		return nil, errors.BadRequestError("Appointment time conflicts with existing appointments", fmt.Errorf("appointment time conflicts with existing appointments"))
	}
	appointment, updateError := queries.UpdateAppointment(ctx, gen_queries.UpdateAppointmentParams{
		PatientID: uuid.NullUUID{},
		DoctorID:  uuid.NullUUID{},
		StartTime: newStartTime,
		EndTime:   newEndTime,
		ID:        uuid.MustParse(appointmentID),
	})
	if updateError != nil {
		return nil, updateError
	}
	return &domain.Appointment{
		ID:        appointment.ID.String(),
		PatientID: appointment.PatientID.UUID.String(),
		DoctorID:  appointment.DoctorID.UUID.String(),
		StartTime: appointment.StartTime,
		EndTime:   appointment.EndTime,
		Status:    appointment.Status,
	}, nil
}

func (r *appointmentRepository) CancelAppointment(appointmentID string, patientID string, reason string) error {
	ctx := context.Background()
	queries := gen_queries.New(r.db)
	_, err := queries.AddAppointmentCancellation(ctx, gen_queries.AddAppointmentCancellationParams{
		AppointmentID: uuid.MustParse(appointmentID),
		PatientID:     uuid.MustParse(patientID),
		Reason:        reason,
	})
	if err != nil {
		return err
	}
	err = queries.UpdateAppointmentStatus(ctx, gen_queries.UpdateAppointmentStatusParams{
		Status: "canceled",
		ID:     uuid.MustParse(appointmentID),
	})
	return err
}

func (r *appointmentRepository) GetAppointmentsByDoctorID(doctorID string) ([]*domain.Appointment, error) {
	ctx := context.Background()
	queries := gen_queries.New(r.db)

	query, err := queries.GetAppointmentsByDoctorID(ctx, uuid.NullUUID{
		UUID:  uuid.MustParse(doctorID),
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	var appointments []*domain.Appointment
	for _, q := range query {
		appointment := &domain.Appointment{
			ID:        q.ID.String(),
			PatientID: q.PatientID.UUID.String(),
			DoctorID:  q.DoctorID.UUID.String(),
			StartTime: q.StartTime,
			EndTime:   q.EndTime,
			Status:    q.Status,
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil

}

func (r *appointmentRepository) GetAppointmentsByDoctorIDAndStartAndEndTime(doctorID string, startTime time.Time, endTime time.Time) ([]*domain.Appointment, error) {
	ctx := context.Background()

	queries := gen_queries.New(r.db)
	query, err := queries.GetAppointmentsByDoctorIDAndStartAndEndTime(ctx, gen_queries.GetAppointmentsByDoctorIDAndStartAndEndTimeParams{
		DoctorID:  uuid.NullUUID{UUID: uuid.MustParse(doctorID), Valid: true},
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		return nil, err
	}

	var appointments []*domain.Appointment
	for _, q := range query {
		appointment := &domain.Appointment{
			ID:        q.ID.String(),
			PatientID: q.PatientID.UUID.String(),
			DoctorID:  q.DoctorID.UUID.String(),
			StartTime: q.StartTime,
			EndTime:   q.EndTime,
			Status:    q.Status,
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *appointmentRepository) GetAppointmentsByPatientID(patientID string) ([]*domain.Appointment, error) {
	ctx := context.Background()
	queries := gen_queries.New(r.db)
	appointmentsQuery, err := queries.GetUpcomingAppointmentsByPatientID(ctx, uuid.NullUUID{UUID: uuid.MustParse(patientID), Valid: true})
	if err != nil {
		return nil, err
	}

	var appointments []*domain.Appointment
	for _, q := range appointmentsQuery {
		appointment := &domain.Appointment{
			ID:        q.ID.String(),
			PatientID: q.PatientID.UUID.String(),
			DoctorID:  q.DoctorID.UUID.String(),
			StartTime: q.StartTime,
			EndTime:   q.EndTime,
			Status:    q.Status,
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil

}

func getTimeUntilNonUtc(startTime time.Time) time.Duration {
	// Get the current time in UTC
	currentTimeUTC := time.Now()

	startTimeUTC := startTime

	// Calculate the duration until the start time
	durationUntilStart := startTimeUTC.Sub(currentTimeUTC)

	return durationUntilStart
}
