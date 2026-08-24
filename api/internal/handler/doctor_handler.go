package handler

import (
	"github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/service"
	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	svc service.DoctorService
}

func NewDoctorHandler(svc service.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		svc: svc,
	}
}

func (h *DoctorHandler) GetDoctors(ctx *gin.Context) {

	doctors, err := h.svc.GetDoctors()
	if err != nil {
		ctx.JSON(500, errors.InternalServerError(err))
		return
	}
	ctx.JSON(200, doctors)

}
