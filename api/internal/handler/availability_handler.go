package handler

import (
	"time"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/errors"
	"github.com/SaidiBTW/appointment_booking_system_go/internal/service"

	"github.com/gin-gonic/gin"
)

type AvailabilityHandler struct {
	svc service.AvailabilityService
}

func NewAvailabilityHandler(svc service.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{
		svc: svc,
	}
}

func (h *AvailabilityHandler) GetAvailability(ctx *gin.Context) {
	Id := ctx.Param("id")
	dateStr := ctx.Query("date")
	Date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(400, errors.BadRequestError("Invalid date format. Use YYYY-MM-DD", err))
		return
	}
	// ctx.Error(errors.NotFoundError("Doctor not found"))

	// return
	slots, err := h.svc.GetAvailableSlots(Id, Date)
	if err != nil {
		ctx.JSON(500, errors.InternalServerError(err))
		return
	}

	ctx.JSON(200, gin.H{"available_slots": slots})
}
