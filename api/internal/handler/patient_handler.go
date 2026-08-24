package handler

import (
	"github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/service"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	svc service.PatientService
}

func NewPatientHandler(svc service.PatientService) *PatientHandler {
	return &PatientHandler{
		svc: svc,
	}
}

func (h *PatientHandler) GetPatients(ctx *gin.Context) {
	patients, err := h.svc.GetPatients()
	if err != nil {
		ctx.JSON(500, errors.InternalServerError(err))
		return
	}
	ctx.JSON(200, patients)
}
