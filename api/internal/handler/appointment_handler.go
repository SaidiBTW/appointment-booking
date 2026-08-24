package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/dto"
	"github.com/go-playground/validator/v10"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/service"
	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	svc service.AppointmentService
}

var Vaidate = validator.New()

func NewAppointmentHandler(svc service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		svc: svc,
	}
}

func (h *AppointmentHandler) CreateAppointment(ctx *gin.Context) {
	var appointmentDto dto.CreateAppointmentRequest

	if err := ctx.ShouldBindJSON(&appointmentDto); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	err := Vaidate.Struct(appointmentDto)

	startTime := appointmentDto.StartTime.Truncate(time.Minute * 15)
	endTime := appointmentDto.EndTime.Truncate(time.Minute * 15)

	if startTime.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("start time cannot be in the past", fmt.Errorf("start time cannot be in the past")))
		return
	}

	if endTime.Sub(startTime) <= 0 {
		// return nil, errors.BadRequestError("end time must be after start time", fmt.Errorf("end time must be after start time"))
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("end time must be after start time", fmt.Errorf("end time must be after start time")))
		return
	}
	if endTime.Sub(startTime) != 30*time.Minute {
		ctx.JSON(http.StatusBadRequest, errors.CustomError(400, "bad_request", "end time must be exactly 30 minutes after start time"))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("Invalid request body", fmt.Errorf("invalid request body: %v", err.Error())))
		return
	}

	if time.Until(startTime) < time.Hour {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("Cannot create an appointment within 1 hour of the start time", fmt.Errorf("cannot create an appointment within 1 hour of the start time")))
		return
	}

	appointmentDto.StartTime = startTime
	appointmentDto.EndTime = endTime

	res, err := h.svc.CreateAppointment(appointmentDto)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *AppointmentHandler) CancelAppointment(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	patientID := ctx.Query("patient_id")
	reason := ctx.Query("reason")

	if appointmentID == "" || patientID == "" || reason == "" {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("Missing required parameters", nil))
		return
	}

	err := h.svc.CancelAppointment(appointmentID, patientID, reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Appointment canceled successfully",
	})
}

func (h *AppointmentHandler) RescheduleAppointment(ctx *gin.Context) {

	appointmentID := ctx.Param("id")
	newStartTime := ctx.Query("new_start_time")
	newEndTime := ctx.Query("new_end_time")

	if appointmentID == "" || newStartTime == "" || newEndTime == "" {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("Missing required parameters", nil))
		return
	}

	startTimeParsed, err := time.Parse(time.RFC3339, newStartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError(fmt.Sprintf("Invalid date format. Please use RFC3339 format. Start time: %s", newStartTime), err))
		return
	}

	endTimeParsed, err := time.Parse(time.RFC3339, newEndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError(fmt.Sprintf("Invalid date format. Please use RFC3339 format. End time: %s", newEndTime), err))
		return
	}

	if startTimeParsed.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("Start time cannot be in the past", fmt.Errorf("start time cannot be in the past")))
		return
	}

	if endTimeParsed.Sub(startTimeParsed) <= 0 {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("End time must be after start time", fmt.Errorf("end time must be after start time")))
		return
	}

	if endTimeParsed.Sub(startTimeParsed) != 30*time.Minute {
		ctx.JSON(http.StatusBadRequest, errors.CustomError(400, "bad_request", "End time must be exactly 30 minutes after start time"))
		return
	}
	if time.Until(startTimeParsed) < time.Hour {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("Cannot reschedule an appointment within 1 hour of the start time", fmt.Errorf("cannot reschedule an appointment within 1 hour of the start time")))
		return
	}

	appointment, rescheduleError := h.svc.RescheduleAppointment(appointmentID, startTimeParsed, endTimeParsed)
	if rescheduleError != nil {
		ctx.JSON(http.StatusInternalServerError, errors.InternalServerError(rescheduleError))
		return
	}

	ctx.JSON(http.StatusOK, appointment)

}

func (h *AppointmentHandler) GetAppointmentsByPatientID(ctx *gin.Context) {
	patientID := ctx.Param("id")

	if patientID == "" {
		ctx.JSON(http.StatusBadRequest, errors.BadRequestError("Missing required parameter: patient_id", nil))
		return
	}

	appointments, err := h.svc.GetAppointmentsByPatientID(patientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		return
	}

	ctx.JSON(http.StatusOK, appointments)
}
