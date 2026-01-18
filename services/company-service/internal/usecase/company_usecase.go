package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vipos89/timehub/services/company-service/internal/domain"
)

type companyUsecase struct {
	repo           domain.CompanyRepository
	contextTimeout time.Duration
	authServiceURL string
}

func NewCompanyUsecase(repo domain.CompanyRepository, timeout time.Duration, authServiceURL string) domain.CompanyUsecase {
	return &companyUsecase{
		repo:           repo,
		contextTimeout: timeout,
		authServiceURL: authServiceURL,
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

func (u *companyUsecase) GetMyCompanies(ctx context.Context, ownerID uint) ([]domain.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetCompaniesByOwnerID(ctx, ownerID)
}

func (u *companyUsecase) GetCompanyByID(ctx context.Context, id uint) (*domain.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetCompanyByID(ctx, id)
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

func (u *companyUsecase) GetCompanyBranches(ctx context.Context, companyID uint) ([]domain.Branch, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetBranchesByCompanyID(ctx, companyID)
}

func (u *companyUsecase) AddCategory(ctx context.Context, companyID uint, branchID uint, name string) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	category := &domain.Category{
		CompanyID: companyID,
		BranchID:  branchID,
		Name:      name,
	}
	err := u.repo.CreateCategory(ctx, category)
	return category, err
}

func (u *companyUsecase) GetBranchCategories(ctx context.Context, branchID uint) ([]domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetCategoriesByBranchID(ctx, branchID)
}

func (u *companyUsecase) AddService(ctx context.Context, companyID uint, branchID uint, categoryID *uint, name, description string, price float64, duration int) (*domain.Service, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	service := &domain.Service{
		CompanyID:       companyID,
		BranchID:        branchID,
		CategoryID:      categoryID,
		Name:            name,
		Description:     description,
		Price:           price,
		DurationMinutes: duration,
	}
	err := u.repo.CreateService(ctx, service)
	return service, err
}

func (u *companyUsecase) UpdateService(ctx context.Context, service *domain.Service) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.UpdateService(ctx, service)
}

func (u *companyUsecase) GetBranchServices(ctx context.Context, branchID uint) ([]domain.Service, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetServicesByBranchID(ctx, branchID)
}

func (u *companyUsecase) AddEmployee(ctx context.Context, branchID uint, name, position, email string) (*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	var userID *uint
	if email != "" {
		// Call Auth Service to create user
		// In a real app, this should be gRPC or more robust HTTP client
		// Using simple http.Post for brevity as this is a demo/prototype phase
		userData := map[string]string{
			"email":    email,
			"password": "temporary_password_123", // Should be random or sent via email
			"role":     "master",
		}
		jsonBody, _ := json.Marshal(userData)
		resp, err := http.Post(u.authServiceURL+"/auth/register", "application/json", bytes.NewBuffer(jsonBody))
		if err == nil && resp.StatusCode == http.StatusCreated {
			var result struct {
				ID uint `json:"id"`
			}
			json.NewDecoder(resp.Body).Decode(&result)
			userID = &result.ID
			resp.Body.Close()
		}
	}

	employee := &domain.Employee{
		BranchID: branchID,
		UserID:   userID,
		Name:     name,
		Position: position,
	}
	err := u.repo.CreateEmployee(ctx, employee)
	return employee, err
}

func (u *companyUsecase) GetCompanyEmployees(ctx context.Context, companyID uint) ([]domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetEmployeesByCompanyID(ctx, companyID)
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

func (u *companyUsecase) RemoveService(ctx context.Context, employeeID, serviceID uint) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.RemoveServiceFromEmployee(ctx, employeeID, serviceID)
}

func (u *companyUsecase) GetEmployeeMenu(ctx context.Context, employeeID uint) ([]domain.EmployeeService, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetServicesByEmployeeID(ctx, employeeID)
}
