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

func (r *companyRepository) UpdateService(ctx context.Context, service *domain.Service) error {
	return r.db.WithContext(ctx).Model(&domain.Service{ID: service.ID}).Updates(service).Error
}

func (r *companyRepository) CreateEmployee(ctx context.Context, employee *domain.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

func (r *companyRepository) GetCompanyByID(ctx context.Context, id uint) (*domain.Company, error) {
	var company domain.Company
	err := r.db.WithContext(ctx).
		Preload("Branches.Employees").
		Preload("Services").
		First(&company, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetCompaniesByOwnerID(ctx context.Context, ownerID uint) ([]domain.Company, error) {
	var companies []domain.Company
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&companies).Error
	return companies, err
}

func (r *companyRepository) GetBranchesByCompanyID(ctx context.Context, companyID uint) ([]domain.Branch, error) {
	var branches []domain.Branch
	err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&branches).Error
	return branches, err
}

func (r *companyRepository) GetCategoriesByBranchID(ctx context.Context, branchID uint) ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.WithContext(ctx).
		Preload("Services").
		Where("branch_id = ?", branchID).
		Find(&categories).Error
	return categories, err
}

func (r *companyRepository) GetServicesByBranchID(ctx context.Context, branchID uint) ([]domain.Service, error) {
	var services []domain.Service
	err := r.db.WithContext(ctx).
		Where("branch_id = ?", branchID).
		Find(&services).Error
	return services, err
}

func (r *companyRepository) GetEmployeesByCompanyID(ctx context.Context, companyID uint) ([]domain.Employee, error) {
	var employees []domain.Employee
	// Find all employees belonging to any branch of the company
	err := r.db.WithContext(ctx).
		Joins("JOIN branches ON branches.id = employees.branch_id").
		Where("branches.company_id = ?", companyID).
		Preload("Services.Service").
		Find(&employees).Error
	return employees, err
}

func (r *companyRepository) AssignServiceToEmployee(ctx context.Context, relation *domain.EmployeeService) error {
	// Upsert: On conflict update price/duration
	return r.db.WithContext(ctx).Save(relation).Error
}

func (r *companyRepository) RemoveServiceFromEmployee(ctx context.Context, employeeID, serviceID uint) error {
	return r.db.WithContext(ctx).Delete(&domain.EmployeeService{}, "employee_id = ? AND service_id = ?", employeeID, serviceID).Error
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
