package repository

import (
	"context"
	"database/sql"

	"log"
	"time"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository/gen_queries"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/shared"
	"github.com/google/uuid"
)

type AppointmentSlot struct {
	StartTime shared.TimeOnly
	EndTime   shared.TimeOnly
}

func (s AppointmentSlot) String() string {
	return s.StartTime.String() + " - " + s.EndTime.String()
}

type AvailabilityRepository interface {
	GetAvailableSlots(doctorID string, date time.Time) ([]AppointmentSlot, error)
}

type availabilityRepository struct {
	db *sql.DB
}

func NewPostgresAvailabilityRepository(db *sql.DB) AvailabilityRepository {
	return &availabilityRepository{
		db: db,
	}
}

func (r *availabilityRepository) GetAvailableSlots(doctorID string, date time.Time) ([]AppointmentSlot, error) {
	ctx := context.Background()

	// availableSet := make(map[AppointmentSlot]struct{})
	dayOfWeek := int(date.Weekday())
	queries := gen_queries.New(r.db)
	query, err := queries.GetDoctorSceduleforDayOfWeek(ctx, gen_queries.GetDoctorSceduleforDayOfWeekParams{
		DoctorID:  uuid.MustParse(doctorID),
		DayOfWeek: int16(dayOfWeek),
	})

	if err != nil {
		return nil, err
	}

	var slots []AppointmentSlot
	for _, q := range query {
		slots = append(slots, generateTimeSlots(q.StartTime, q.EndTime, 30*time.Minute)...)
	}

	appointments, err := r.GetAppointmentsByDoctorID(doctorID, date)
	if err != nil {
		return nil, err
	}
	slots = removeBookedSlots(slots, appointments)

	return slots, nil
}

func removeBookedSlots(availableSlots []AppointmentSlot, bookedSlots []AppointmentSlot) []AppointmentSlot {
	bookedSet := make(map[string]struct{})
	for _, booked := range bookedSlots {
		bookedSet[booked.String()] = struct{}{}
	}

	var filteredSlots []AppointmentSlot
	for _, available := range availableSlots {
		if _, exists := bookedSet[available.String()]; !exists {
			filteredSlots = append(filteredSlots, available)
		}
	}

	return filteredSlots
}

func (r *availabilityRepository) GetAppointmentsByDoctorID(doctorID string, date time.Time) ([]AppointmentSlot, error) {
	ctx := context.Background()

	queries := gen_queries.New(r.db)
	dayOfWeek := int(date.Weekday())
	log.Printf("Getting appointments for doctorID: %s on date: %s (day of week: %d)", doctorID, date.Format("2006-01-02"), dayOfWeek)
	appointments, err := queries.GetAppointmentsByDoctorID(ctx, uuid.NullUUID{
		UUID:  uuid.MustParse(doctorID),
		Valid: true,
	})

	if err != nil {
		return nil, err
	}
	var appointmentSlots []AppointmentSlot
	for _, appointment := range appointments {
		if appointment.StartTime.Truncate(24 * time.Hour).Equal(date.Truncate(24 * time.Hour)) {
			appointmentSlots = append(appointmentSlots, AppointmentSlot{
				StartTime: shared.TimeOnly(appointment.StartTime),
				EndTime:   shared.TimeOnly(appointment.EndTime),
			})
		}
	}

	return appointmentSlots, nil
}

func generateTimeSlots(startTime, endTime time.Time, slotDuration time.Duration) []AppointmentSlot {
	var slots []AppointmentSlot
	for current := startTime; current.Before(endTime) || current.Equal(endTime); current = current.Add(slotDuration) {
		slots = append(slots, AppointmentSlot{
			StartTime: shared.TimeOnly(current),
			EndTime:   shared.TimeOnly(current.Add(slotDuration)),
		})
	}
	return slots
}
