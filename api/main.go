package main

import (
	"fmt"
	"log"
	"os"

	_ "time/tzdata"

	"github.com/SaidiBTW/appointment_booking_system_go/config"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/handler"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/repository"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/service"
	"github.com/SaidiBTW/appointment_booking_system_go/middleware"
	"github.com/SaidiBTW/appointment_booking_system_go/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	postgresConfig, err := config.Load()
	fmt.Println("Postgres Config:", postgresConfig)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := storage.NewPostgresDB(postgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	apointmentRepo := repository.NewPostgresAppointmentRepository(db)

	appointmentService := service.NewAppointmentService(apointmentRepo)

	appointmentHandler := handler.NewAppointmentHandler(appointmentService)

	availabilityRepo := repository.NewPostgresAvailabilityRepository(db)

	availabilityService := service.NewAvailabilityService(availabilityRepo)

	availabilityHandler := handler.NewAvailabilityHandler(availabilityService)

	patientRepo := repository.NewPostgresPatientRepository(db)

	patientService := service.NewPatientService(patientRepo)

	patientHandler := handler.NewPatientHandler(patientService)

	doctorRepo := repository.NewPostgresDoctorRepository(db)

	doctorService := service.NewDoctorService(doctorRepo)

	doctorHandler := handler.NewDoctorHandler(doctorService)

	router := gin.Default()

	router.Use(middleware.ErrorHandler())

	api := router.Group("/api/v1")
	{
		api.POST("/appointments", appointmentHandler.CreateAppointment)
		api.PATCH("/appointments/:id/cancel", appointmentHandler.CancelAppointment)
		api.PATCH("/appointments/:id/reschedule", appointmentHandler.RescheduleAppointment)
	}
	{
		api.GET("/doctors/:id/availability", availabilityHandler.GetAvailability)
		api.GET("/doctors", doctorHandler.GetDoctors)
	}
	{
		api.GET("/patients/:id/appointments", appointmentHandler.GetAppointmentsByPatientID)
		api.GET("/patients", patientHandler.GetPatients)
	}
	{
		api.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"status": "ok",
			})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified in .env
	}

	log.Printf("Starting server on port %s", port)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
