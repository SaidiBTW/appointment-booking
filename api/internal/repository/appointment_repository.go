package repository

import (
	"context"
	"database/sql"
	"log"

	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"errors"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/domain"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/dto"
	custom_errors "github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository/gen_queries"
)

type AppointmentRepository interface {
	CreateAppointment(appointmentDto dto.CreateAppointmentRequest) (*domain.Appointment, error)
	GetAppointmentByID(id string) (*domain.Appointment, error)
	GetAppointmentsByDoctorID(doctorID string) ([]*domain.Appointment, error)
	GetAppointmentsByPatientID(patientID string) ([]*domain.Appointment, error)
	RescheduleAppointment(appointmentID string, newStartTime time.Time, newEndTime time.Time) (*domain.Appointment, error)
	CancelAppointment(appointmentID string, patientID string, reason string) (*domain.AppointmentCancellation, error)
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
		// Check whether error is due to a unique constraint violation (appointment already exists)
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, custom_errors.BadRequestError("Appointment already exists", err)
		}

		return nil, custom_errors.InternalServerError(err)
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

	appointment, err := r.GetAppointmentByID(appointmentID)
	if err != nil {
		return nil, err
	}

	if appointment.Status == "canceled" {
		return nil, custom_errors.BadRequestError("Cannot reschedule a canceled appointment", fmt.Errorf("cannot reschedule a canceled appointment"))
	}

	appointmentsInRange, err := r.GetAppointmentsByDoctorIDAndStartAndEndTime(appointment.DoctorID, newStartTime, newEndTime)
	if err != nil {
		return nil, err
	}
	if len(appointmentsInRange) > 0 {
		return nil, custom_errors.BadRequestError("Appointment time conflicts with existing appointments", fmt.Errorf("appointment time conflicts with existing appointments"))
	}
	updatedAppointment, updateError := queries.UpdateAppointment(ctx, gen_queries.UpdateAppointmentParams{
		PatientID: uuid.NullUUID{},
		DoctorID:  uuid.NullUUID{},
		StartTime: newStartTime,
		EndTime:   newEndTime,
		ID:        uuid.MustParse(appointmentID),
		Status:    "scheduled",
	})
	if updateError != nil {
		return nil, updateError
	}
	return &domain.Appointment{
		ID:        updatedAppointment.ID.String(),
		PatientID: updatedAppointment.PatientID.UUID.String(),
		DoctorID:  updatedAppointment.DoctorID.UUID.String(),
		StartTime: updatedAppointment.StartTime,
		EndTime:   updatedAppointment.EndTime,
		Status:    updatedAppointment.Status,
	}, nil
}

func (r *appointmentRepository) CancelAppointment(appointmentID string, patientID string, reason string) (*domain.AppointmentCancellation, error) {
	ctx := context.Background()
	queries := gen_queries.New(r.db)
	cancellation, err := queries.AddAppointmentCancellation(ctx, gen_queries.AddAppointmentCancellationParams{
		AppointmentID: uuid.MustParse(appointmentID),
		PatientID:     uuid.MustParse(patientID),
		Reason:        reason,
	})
	var psqlErr *pq.Error
	if errors.As(err, &psqlErr) && psqlErr.Code == "23505" {
		return nil, custom_errors.BadRequestError("Appointment cancellation already exists", err)
	}
	if err != nil {
		return nil, custom_errors.InternalServerError(err)
	}
	err = queries.UpdateAppointmentStatus(ctx, gen_queries.UpdateAppointmentStatusParams{
		Status: "canceled",
		ID:     uuid.MustParse(appointmentID),
	})
	if err != nil {
		return nil, custom_errors.InternalServerError(err)
	}
	return &domain.AppointmentCancellation{
		ID:            cancellation.ID.String(),
		AppointmentID: cancellation.AppointmentID.String(),
		PatientID:     cancellation.PatientID.String(),
		Reason:        cancellation.Reason,
	}, nil
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
