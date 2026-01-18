package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Company represents a business entity (e.g., a salon chain).
type Company struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	TaxID     string         `json:"tax_id"`
	OwnerID   uint           `json:"owner_id" gorm:"not null"` // Link to Auth Service User
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Branches   []Branch   `json:"branches,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Services   []Service  `json:"services,omitempty"`
}

// Branch represents a physical location of the company.
type Branch struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CompanyID uint           `json:"company_id" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	Address   string         `json:"address"`
	Phone     string         `json:"phone"`
	IsMain    bool           `json:"is_main" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Employees []Employee `json:"employees,omitempty"`
}

// Category groups services (e.g., "Hair", "Nails").
type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CompanyID uint           `json:"company_id" gorm:"not null"`
	BranchID  uint           `json:"branch_id" gorm:"not null;index"`
	Name      string         `json:"name" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Services []Service `json:"services,omitempty" gorm:"foreignKey:CategoryID"`
}

// Service represents a catalog item (definition only, no price/duration).
type Service struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	CompanyID       uint           `json:"company_id" gorm:"not null"`
	BranchID        uint           `json:"branch_id" gorm:"not null;index"`
	CategoryID      *uint          `json:"category_id"` // Optional
	Name            string         `json:"name" gorm:"not null"`
	Description     string         `json:"description"`
	Price           float64        `json:"price" gorm:"not null;default:0"`
	DurationMinutes int            `json:"duration_minutes" gorm:"not null;default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// Employee represents a staff member (Master).
type Employee struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	BranchID  uint           `json:"branch_id" gorm:"not null"`
	UserID    *uint          `json:"user_id"` // Optional link to Auth Service
	Name      string         `json:"name" gorm:"not null"`
	Position  string         `json:"position"`
	AvatarURL string         `json:"avatar_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Services []EmployeeService `json:"services,omitempty" gorm:"foreignKey:EmployeeID"`
}

// EmployeeService is the pricing matrix (Junction Table).
// Defines that a specific Employee performs a specific Service at a specific Price/Duration.
type EmployeeService struct {
	EmployeeID      uint    `json:"employee_id" gorm:"primaryKey"`
	ServiceID       uint    `json:"service_id" gorm:"primaryKey"`
	Price           float64 `json:"price" gorm:"not null"`
	DurationMinutes int     `json:"duration_minutes" gorm:"not null"`

	// Relations for Preloading
	Service *Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

// Interfaces

type CompanyRepository interface {
	// Transactions
	CreateCompanyWithBranch(ctx context.Context, company *Company, branch *Branch) error

	// CRUD
	CreateBranch(ctx context.Context, branch *Branch) error
	CreateCategory(ctx context.Context, category *Category) error
	CreateService(ctx context.Context, service *Service) error
	UpdateService(ctx context.Context, service *Service) error
	CreateEmployee(ctx context.Context, employee *Employee) error

	// Queries
	GetCompanyByID(ctx context.Context, id uint) (*Company, error)
	GetCompaniesByOwnerID(ctx context.Context, ownerID uint) ([]Company, error)
	GetBranchesByCompanyID(ctx context.Context, companyID uint) ([]Branch, error)
	GetCategoriesByBranchID(ctx context.Context, branchID uint) ([]Category, error)
	GetServicesByBranchID(ctx context.Context, branchID uint) ([]Service, error)
	GetEmployeesByCompanyID(ctx context.Context, companyID uint) ([]Employee, error)
	GetServicesByEmployeeID(ctx context.Context, employeeID uint) ([]EmployeeService, error)

	// Pricing Matrix
	AssignServiceToEmployee(ctx context.Context, relation *EmployeeService) error
	RemoveServiceFromEmployee(ctx context.Context, employeeID, serviceID uint) error
}

type CompanyUsecase interface {
	CreateCompany(ctx context.Context, name string, ownerID uint) (*Company, error)
	GetMyCompanies(ctx context.Context, ownerID uint) ([]Company, error)
	GetCompanyByID(ctx context.Context, id uint) (*Company, error)

	AddBranch(ctx context.Context, companyID uint, name, address, phone string) (*Branch, error)
	GetCompanyBranches(ctx context.Context, companyID uint) ([]Branch, error)

	AddCategory(ctx context.Context, companyID uint, branchID uint, name string) (*Category, error)
	GetBranchCategories(ctx context.Context, branchID uint) ([]Category, error)

	AddService(ctx context.Context, companyID uint, branchID uint, categoryID *uint, name, description string, price float64, duration int) (*Service, error)
	UpdateService(ctx context.Context, service *Service) error
	GetBranchServices(ctx context.Context, branchID uint) ([]Service, error)

	AddEmployee(ctx context.Context, branchID uint, name, position, email string) (*Employee, error)
	GetCompanyEmployees(ctx context.Context, companyID uint) ([]Employee, error)

	AssignService(ctx context.Context, employeeID, serviceID uint, price float64, duration int) error
	RemoveService(ctx context.Context, employeeID, serviceID uint) error
	GetEmployeeMenu(ctx context.Context, employeeID uint) ([]EmployeeService, error)
}
