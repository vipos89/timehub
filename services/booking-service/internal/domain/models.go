package domain

import (
	"context"
	"time"
)

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "pending"
	StatusConfirmed AppointmentStatus = "confirmed"
	StatusCancelled AppointmentStatus = "cancelled"
)

type Schedule struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	EmployeeID uint      `json:"employee_id" gorm:"not null;index"`
	DayOfWeek  int       `json:"day_of_week" gorm:"not null"` // 0 = Sunday, 1 = Monday, ...
	StartTime  string    `json:"start_time" gorm:"not null"`  // e.g. "09:00"
	EndTime    string    `json:"end_time" gorm:"not null"`    // e.g. "18:00"
	IsDayOff   bool      `json:"is_day_off" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type WorkShift struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	EmployeeID uint      `json:"employee_id" gorm:"not null;index"`
	BranchID   uint      `json:"branch_id" gorm:"not null;index"`
	Date       time.Time `json:"date" gorm:"type:date;not null;index"`
	StartTime  string    `json:"start_time" gorm:"not null"` // e.g. "09:00"
	EndTime    string    `json:"end_time" gorm:"not null"`   // e.g. "18:00"
	IsDayOff   bool      `json:"is_day_off" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Appointment struct {
	ID         uint              `json:"id" gorm:"primaryKey"`
	EmployeeID uint              `json:"employee_id" gorm:"not null;index"`
	ServiceID  uint              `json:"service_id" gorm:"not null"`
	ClientID   uint              `json:"client_id" gorm:"not null;index"`
	StartTime  time.Time         `json:"start_time" gorm:"not null;index"`
	EndTime    time.Time         `json:"end_time" gorm:"not null"`
	Status     AppointmentStatus `json:"status" gorm:"type:text;default:'pending'"`
	Comment    string            `json:"comment"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type Slot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsFree    bool      `json:"is_free"`
}

type BookingRepository interface {
	// Schedule
	GetScheduleByEmployee(ctx context.Context, employeeID uint) ([]Schedule, error)
	UpdateSchedule(ctx context.Context, schedule []Schedule) error

	// Appointments
	CreateAppointment(ctx context.Context, appointment *Appointment) error
	GetAppointmentsByEmployee(ctx context.Context, employeeID uint, start, end time.Time) ([]Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, id uint, status AppointmentStatus) error

	// Work Shifts
	GetShiftsByEmployee(ctx context.Context, employeeID uint, start, end time.Time) ([]WorkShift, error)
	GetShiftsByBranch(ctx context.Context, branchID uint, start, end time.Time) ([]WorkShift, error)
	UpsertShifts(ctx context.Context, shifts []WorkShift) error
}

type BookingUsecase interface {
	GetAvailableSlots(ctx context.Context, employeeID uint, serviceID uint, date time.Time) ([]Slot, error)
	CreateBooking(ctx context.Context, appointment *Appointment) error
	GetEmployeeSchedule(ctx context.Context, employeeID uint) ([]Schedule, error)
	SetEmployeeSchedule(ctx context.Context, employeeID uint, schedules []Schedule) error

	// Work Shifts
	GetShifts(ctx context.Context, employeeID uint, branchID uint, month time.Time) ([]WorkShift, error)
	SaveShifts(ctx context.Context, shifts []WorkShift) error
}
