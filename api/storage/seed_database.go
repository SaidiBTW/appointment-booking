package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository/gen_queries"
)

func SeedDatabase(db *sql.DB) {
	ctx := context.Background()
	// Create a new instance of the generated queries
	queries := gen_queries.New(db)

	// Seed the database with initial data
	// Example: Create a doctor
	doctor1, _ := queries.CreateDoctor(ctx, "doctor-1")
	doctor2, _ := queries.CreateDoctor(ctx, "doctor-2")
	doctor3, _ := queries.CreateDoctor(ctx, "doctor-3")
	doctor4, _ := queries.CreateDoctor(ctx, "doctor-4")
	doctor5, _ := queries.CreateDoctor(ctx, "doctor-5")
	doctor6, _ := queries.CreateDoctor(ctx, "doctor-6")
	doctor7, _ := queries.CreateDoctor(ctx, "doctor-7")
	doctor8, _ := queries.CreateDoctor(ctx, "doctor-8")
	doctor9, _ := queries.CreateDoctor(ctx, "doctor-9")
	doctor10, _ := queries.CreateDoctor(ctx, "doctor-10")

	// Create schedules for each doctor
	// Example: Create a schedule for doctor-1 on Monday from 9 AM to 5 PM
	for _, doctor := range []gen_queries.Doctor{doctor1, doctor2, doctor3, doctor4, doctor5, doctor6, doctor7, doctor8, doctor9, doctor10} {
		for dayOfWeek := 1; dayOfWeek <= 5; dayOfWeek++ { // Monday to Friday
			queries.CreateSchedule(ctx, gen_queries.CreateScheduleParams{
				DoctorID:  doctor.ID,
				DayOfWeek: int16(dayOfWeek),
				StartTime: (func() time.Time { n, _ := time.Parse(time.TimeOnly, "09:00:00"); return n })(),
				EndTime:   (func() time.Time { n, _ := time.Parse(time.TimeOnly, "17:00:00"); return n })(),
			})
		}
	}

	// Example: Create a patient
	for i := 1; i <= 10; i++ {
		queries.CreatePatient(ctx, gen_queries.CreatePatientParams{
			Name:  "patient-" + fmt.Sprint(i),
			Email: "patient" + fmt.Sprint(i) + "@example.com",
		})
	}

}
