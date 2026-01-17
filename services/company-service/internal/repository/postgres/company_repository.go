package postgres

import (
	"context"
	"errors"

	"github.com/vipos89/timehub/services/company-service/internal/domain"
	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) domain.CompanyRepository {
	return &companyRepository{db: db}
}

// CreateCompanyWithBranch implements Transactional creation
func (r *companyRepository) CreateCompanyWithBranch(ctx context.Context, company *domain.Company, branch *domain.Branch) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create Company
		if err := tx.Create(company).Error; err != nil {
			return err
		}

		// 2. Link Branch to Company
		branch.CompanyID = company.ID
		branch.IsMain = true // Enforce main branch on creation

		// 3. Create Branch
		if err := tx.Create(branch).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *companyRepository) CreateBranch(ctx context.Context, branch *domain.Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

func (r *companyRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *companyRepository) CreateService(ctx context.Context, service *domain.Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *companyRepository) CreateEmployee(ctx context.Context, employee *domain.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

func (r *companyRepository) GetCompanyByID(ctx context.Context, id uint) (*domain.Company, error) {
	var company domain.Company
	err := r.db.WithContext(ctx).
		Preload("Branches").
		Preload("Categories").
		Preload("Services").
		First(&company, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error, return nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) AssignServiceToEmployee(ctx context.Context, relation *domain.EmployeeService) error {
	// Upsert: On conflict update price/duration
	return r.db.WithContext(ctx).Save(relation).Error
}

func (r *companyRepository) GetServicesByEmployeeID(ctx context.Context, employeeID uint) ([]domain.EmployeeService, error) {
	var results []domain.EmployeeService
	// Join with Service table to get names
	err := r.db.WithContext(ctx).
		Preload("Service").
		Where("employee_id = ?", employeeID).
		Find(&results).Error
	return results, err
}
