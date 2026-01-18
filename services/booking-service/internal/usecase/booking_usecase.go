package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/vipos89/timehub/services/booking-service/internal/domain"
)

type bookingUsecase struct {
	repo    domain.BookingRepository
	timeout time.Duration
}

func NewBookingUsecase(repo domain.BookingRepository, timeout time.Duration) domain.BookingUsecase {
	return &bookingUsecase{
		repo:    repo,
		timeout: timeout,
	}
}

func (u *bookingUsecase) GetAvailableSlots(ctx context.Context, employeeID uint, serviceID uint, date time.Time) ([]domain.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// 1. Check for WorkShift override (date-specific)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)

	shifts, err := u.repo.GetShiftsByEmployee(ctx, employeeID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	var startTime, endTime string
	var isDayOff bool
	var hasOverride bool

	if len(shifts) > 0 {
		startTime = shifts[0].StartTime
		endTime = shifts[0].EndTime
		isDayOff = shifts[0].IsDayOff
		hasOverride = true
	} else {
		// 2. Fallback to weekly Schedule template
		dayOfWeek := int(date.Weekday())
		schedules, err := u.repo.GetScheduleByEmployee(ctx, employeeID)
		if err != nil {
			return nil, err
		}

		for _, s := range schedules {
			if s.DayOfWeek == dayOfWeek {
				startTime = s.StartTime
				endTime = s.EndTime
				isDayOff = s.IsDayOff
				hasOverride = true
				break
			}
		}
	}

	if !hasOverride || isDayOff {
		return []domain.Slot{}, nil
	}

	// 3. Define working hours for the date
	startStr := fmt.Sprintf("%s %s:00", date.Format("2006-01-02"), startTime)
	endStr := fmt.Sprintf("%s %s:00", date.Format("2006-01-02"), endTime)

	workingStart, _ := time.ParseInLocation("2006-01-02 15:04:05", startStr, date.Location())
	workingEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", endStr, date.Location())

	// 3. Get existing appointments for this day
	appointments, err := u.repo.GetAppointmentsByEmployee(ctx, employeeID, workingStart, workingEnd)
	if err != nil {
		return nil, err
	}

	// 4. Generate slots (simplified: 30 min duration for now, should get from Service)
	duration := 30 * time.Minute // TODO: Fetch from company service / employee_service price matrix

	var slots []domain.Slot
	for t := workingStart; t.Add(duration).Before(workingEnd) || t.Add(duration).Equal(workingEnd); t = t.Add(duration) {
		slot := domain.Slot{
			StartTime: t,
			EndTime:   t.Add(duration),
			IsFree:    true,
		}

		// Check overlap with appointments
		for _, app := range appointments {
			if t.Before(app.EndTime) && slot.EndTime.After(app.StartTime) {
				slot.IsFree = false
				break
			}
		}
		slots = append(slots, slot)
	}

	return slots, nil
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, appointment *domain.Appointment) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// TODO: Add concurrency check (transactions)
	// 1. Double check availability
	slots, err := u.GetAvailableSlots(ctx, appointment.EmployeeID, appointment.ServiceID, appointment.StartTime)
	if err != nil {
		return err
	}

	available := false
	for _, s := range slots {
		if s.StartTime.Equal(appointment.StartTime) && s.IsFree {
			available = true
			break
		}
	}

	if !available {
		return fmt.Errorf("slot is already taken or out of working hours")
	}

	appointment.Status = domain.StatusConfirmed
	return u.repo.CreateAppointment(ctx, appointment)
}

func (u *bookingUsecase) GetEmployeeSchedule(ctx context.Context, employeeID uint) ([]domain.Schedule, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.repo.GetScheduleByEmployee(ctx, employeeID)
}

func (u *bookingUsecase) SetEmployeeSchedule(ctx context.Context, employeeID uint, schedules []domain.Schedule) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	for i := range schedules {
		schedules[i].EmployeeID = employeeID
	}

	return u.repo.UpdateSchedule(ctx, schedules)
}

func (u *bookingUsecase) GetShifts(ctx context.Context, employeeID uint, branchID uint, month time.Time) ([]domain.WorkShift, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	if branchID > 0 {
		return u.repo.GetShiftsByBranch(ctx, branchID, start, end)
	}
	return u.repo.GetShiftsByEmployee(ctx, employeeID, start, end)
}

func (u *bookingUsecase) SaveShifts(ctx context.Context, shifts []domain.WorkShift) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.repo.UpsertShifts(ctx, shifts)
}
