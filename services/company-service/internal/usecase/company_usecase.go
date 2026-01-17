package usecase

import (
	"context"
	"time"

	"github.com/vipos89/timehub/services/company-service/internal/domain"
)

type companyUsecase struct {
	repo           domain.CompanyRepository
	contextTimeout time.Duration
}

func NewCompanyUsecase(repo domain.CompanyRepository, timeout time.Duration) domain.CompanyUsecase {
	return &companyUsecase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *companyUsecase) CreateCompany(ctx context.Context, name string, ownerID uint) (*domain.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	company := &domain.Company{
		Name:    name,
		OwnerID: ownerID,
	}

	// Create default main branch
	branch := &domain.Branch{
		Name:    "Main Branch",
		Address: "Headquarters",
		IsMain:  true,
	}

	if err := u.repo.CreateCompanyWithBranch(ctx, company, branch); err != nil {
		return nil, err
	}

	return company, nil
}

func (u *companyUsecase) AddBranch(ctx context.Context, companyID uint, name, address, phone string) (*domain.Branch, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	branch := &domain.Branch{
		CompanyID: companyID,
		Name:      name,
		Address:   address,
		Phone:     phone,
		IsMain:    false,
	}

	err := u.repo.CreateBranch(ctx, branch)
	return branch, err
}

func (u *companyUsecase) AddCategory(ctx context.Context, companyID uint, name string) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	category := &domain.Category{
		CompanyID: companyID,
		Name:      name,
	}
	err := u.repo.CreateCategory(ctx, category)
	return category, err
}

func (u *companyUsecase) AddService(ctx context.Context, companyID uint, categoryID *uint, name, description string) (*domain.Service, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	service := &domain.Service{
		CompanyID:   companyID,
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
	}
	err := u.repo.CreateService(ctx, service)
	return service, err
}

func (u *companyUsecase) AddEmployee(ctx context.Context, branchID uint, name, position string) (*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	employee := &domain.Employee{
		BranchID: branchID,
		Name:     name,
		Position: position,
	}
	err := u.repo.CreateEmployee(ctx, employee)
	return employee, err
}

func (u *companyUsecase) AssignService(ctx context.Context, employeeID, serviceID uint, price float64, duration int) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	relation := &domain.EmployeeService{
		EmployeeID:      employeeID,
		ServiceID:       serviceID,
		Price:           price,
		DurationMinutes: duration,
	}
	return u.repo.AssignServiceToEmployee(ctx, relation)
}

func (u *companyUsecase) GetEmployeeMenu(ctx context.Context, employeeID uint) ([]domain.EmployeeService, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetServicesByEmployeeID(ctx, employeeID)
}
