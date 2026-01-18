package http

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vipos89/timehub/pkg/erru"
	"github.com/vipos89/timehub/services/booking-service/internal/domain"
)

type BookingHandler struct {
	Usecase domain.BookingUsecase
}

func NewBookingHandler(e *echo.Echo, us domain.BookingUsecase) {
	handler := &BookingHandler{
		Usecase: us,
	}

	e.GET("/slots", handler.GetSlots)
	e.POST("/bookings", handler.CreateBooking)
	e.GET("/schedules/:employee_id", handler.GetSchedule)
	e.POST("/schedules/:employee_id", handler.SetSchedule)

	// Work Shifts
	e.GET("/shifts", handler.GetShifts)
	e.POST("/shifts", handler.SaveShifts)
}

type getSlotsRequest struct {
	EmployeeID uint      `query:"employee_id" validate:"required"`
	ServiceID  uint      `query:"service_id" validate:"required"`
	Date       time.Time `query:"date" validate:"required" example:"2026-01-20T00:00:00Z"`
}

// GetSlots godoc
// @Summary Get available slots
// @Description Calculate available time slots for an employee and service on a specific date
// @Tags bookings
// @Accept json
// @Produce json
// @Param employee_id query int true "Employee ID"
// @Param service_id query int true "Service ID"
// @Param date query string true "Date (ISO8601)"
// @Success 200 {array} domain.Slot
// @Failure 400 {object} erru.AppError
// @Failure 500 {object} erru.AppError
// @Router /slots [get]
func (h *BookingHandler) GetSlots(c echo.Context) error {
	empID, _ := strconv.Atoi(c.QueryParam("employee_id"))
	svcID, _ := strconv.Atoi(c.QueryParam("service_id"))
	dateStr := c.QueryParam("date")

	date, err := parseDate(dateStr)
	if err != nil {
		log.Printf("Error parsing date '%s': %v", dateStr, err)
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	slots, err := h.Usecase.GetAvailableSlots(c.Request().Context(), uint(empID), uint(svcID), date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, slots)
}

type createBookingRequest struct {
	EmployeeID uint      `json:"employee_id" validate:"required"`
	ServiceID  uint      `json:"service_id" validate:"required"`
	ClientID   uint      `json:"client_id" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required"`
	Comment    string    `json:"comment"`
}

// CreateBooking godoc
// @Summary Book an appointment
// @Description Create a new appointment if the slot is available
// @Tags bookings
// @Accept json
// @Produce json
// @Param body body createBookingRequest true "Booking Info"
// @Success 201 {object} domain.Appointment
// @Failure 400 {object} erru.AppError
// @Failure 500 {object} erru.AppError
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c echo.Context) error {
	var req createBookingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	appointment := &domain.Appointment{
		EmployeeID: req.EmployeeID,
		ServiceID:  req.ServiceID,
		ClientID:   req.ClientID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Comment:    req.Comment,
	}

	err := h.Usecase.CreateBooking(c.Request().Context(), appointment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusCreated, appointment)
}

// GetSchedule godoc
// @Summary Get employee schedule
// @Description Get weekly working schedule for an employee
// @Tags schedules
// @Accept json
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Success 200 {array} domain.Schedule
// @Router /schedules/{employee_id} [get]
func (h *BookingHandler) GetSchedule(c echo.Context) error {
	empID, _ := strconv.Atoi(c.Param("employee_id"))
	schedules, err := h.Usecase.GetEmployeeSchedule(c.Request().Context(), uint(empID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, schedules)
}

// SetSchedule godoc
// @Summary Set employee schedule
// @Description Update working hours for an employee
// @Tags schedules
// @Accept json
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Param body body []domain.Schedule true "Schedules Array"
// @Success 200 {string} string "OK"
// @Router /schedules/{employee_id} [post]
func (h *BookingHandler) SetSchedule(c echo.Context) error {
	empID, _ := strconv.Atoi(c.Param("employee_id"))
	var schedules []domain.Schedule
	if err := c.Bind(&schedules); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	err := h.Usecase.SetEmployeeSchedule(c.Request().Context(), uint(empID), schedules)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, "OK")
}

// GetShifts godoc
// @Summary Get work shifts
// @Description Get shifts for an employee or an entire branch for a specific month
// @Tags shifts
// @Produce json
// @Param employee_id query int false "Employee ID"
// @Param branch_id query int false "Branch ID"
// @Param month query string true "Month (ISO8601, e.g., 2026-01-01T00:00:00Z)"
// @Success 200 {array} domain.WorkShift
// @Router /shifts [get]
func (h *BookingHandler) GetShifts(c echo.Context) error {
	empID, _ := strconv.Atoi(c.QueryParam("employee_id"))
	branchID, _ := strconv.Atoi(c.QueryParam("branch_id"))
	monthStr := c.QueryParam("month")

	month, err := parseDate(monthStr)
	if err != nil {
		log.Printf("Error parsing month '%s': %v", monthStr, err)
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	shifts, err := h.Usecase.GetShifts(c.Request().Context(), uint(empID), uint(branchID), month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, shifts)
}

// SaveShifts godoc
// @Summary Bulk save work shifts
// @Description Create or update multiple work shifts
// @Tags shifts
// @Accept json
// @Produce json
// @Param body body []domain.WorkShift true "Shifts Array"
// @Success 200 {string} string "OK"
// @Router /shifts [post]
func (h *BookingHandler) SaveShifts(c echo.Context) error {
	var shifts []domain.WorkShift
	if err := c.Bind(&shifts); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	err := h.Usecase.SaveShifts(c.Request().Context(), shifts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, "OK")
}

func parseDate(s string) (time.Time, error) {
	// Replace space with + (common issue with URL decoding of +)
	s = strings.Replace(s, " ", "+", 1)

	// RFC3339 covers most ISO8601
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}

	// RFC3339Nano covers fractional seconds
	t, err = time.Parse(time.RFC3339Nano, s)
	if err == nil {
		return t, nil
	}

	// Try without timezone (assume UTC)
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		return t, nil
	}

	// Try with milliseconds but no timezone
	t, err = time.Parse("2006-01-02T15:04:05.000", s)
	if err == nil {
		return t, nil
	}

	// Date only
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		return t, nil
	}

	return time.Parse(time.RFC3339, s) // Final try for meaningful error
}
