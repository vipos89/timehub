package postgres

import (
	"context"
	"time"

	"github.com/vipos89/timehub/services/booking-service/internal/domain"
	"gorm.io/gorm"
)

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) domain.BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) GetScheduleByEmployee(ctx context.Context, employeeID uint) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Order("day_of_week ASC").Find(&schedules).Error
	return schedules, err
}

func (r *bookingRepository) UpdateSchedule(ctx context.Context, schedules []domain.Schedule) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range schedules {
			if err := tx.Save(&s).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *bookingRepository) CreateAppointment(ctx context.Context, appointment *domain.Appointment) error {
	return r.db.WithContext(ctx).Create(appointment).Error
}

func (r *bookingRepository) GetAppointmentsByEmployee(ctx context.Context, employeeID uint, start, end time.Time) ([]domain.Appointment, error) {
	var appointments []domain.Appointment
	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND start_time >= ? AND start_time < ? AND status != ?", employeeID, start, end, domain.StatusCancelled).
		Find(&appointments).Error
	return appointments, err
}

func (r *bookingRepository) UpdateAppointmentStatus(ctx context.Context, id uint, status domain.AppointmentStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Appointment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *bookingRepository) GetShiftsByEmployee(ctx context.Context, employeeID uint, start, end time.Time) ([]domain.WorkShift, error) {
	var shifts []domain.WorkShift
	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, start, end).
		Order("date ASC").
		Find(&shifts).Error
	return shifts, err
}

func (r *bookingRepository) GetShiftsByBranch(ctx context.Context, branchID uint, start, end time.Time) ([]domain.WorkShift, error) {
	var shifts []domain.WorkShift
	err := r.db.WithContext(ctx).
		Where("branch_id = ? AND date >= ? AND date <= ?", branchID, start, end).
		Order("date ASC, employee_id ASC").
		Find(&shifts).Error
	return shifts, err
}

func (r *bookingRepository) UpsertShifts(ctx context.Context, shifts []domain.WorkShift) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range shifts {
			// Upsert logic: find by employee and date
			var existing domain.WorkShift
			res := tx.Where("employee_id = ? AND date = ?", s.EmployeeID, s.Date).First(&existing)
			if res.Error == nil {
				s.ID = existing.ID
				if err := tx.Save(&s).Error; err != nil {
					return err
				}
			} else if res.Error == gorm.ErrRecordNotFound {
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
			} else {
				return res.Error
			}
		}
		return nil
	})
}
