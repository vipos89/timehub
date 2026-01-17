package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/timehub/timehub/pkg/erru"
	"github.com/timehub/timehub/services/company-service/internal/usecase"
)

type CompanyHandler struct {
	Usecase usecase.CompanyUsecase
}

func NewCompanyHandler(e *echo.Echo, us usecase.CompanyUsecase) {
	handler := &CompanyHandler{
		Usecase: us,
	}

	// Company Routes
	e.POST("/companies", handler.CreateCompany)
	e.POST("/companies/:id/branches", handler.AddBranch)
	e.POST("/companies/:id/categories", handler.AddCategory)
	e.POST("/companies/:id/services", handler.AddService)

	// Employee Routes
	e.POST("/branches/:id/employees", handler.AddEmployee)
	e.POST("/employees/:id/services", handler.AssignService)
	e.GET("/employees/:id/services", handler.GetEmployeeMenu)
}

// Request Structs

type createCompanyRequest struct {
	Name string `json:"name" validate:"required"`
}

type addBranchRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type addCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type addServiceRequest struct {
	CategoryID  *uint  `json:"category_id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type addEmployeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Position string `json:"position"`
}

type assignServiceRequest struct {
	ServiceID uint    `json:"service_id" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
	Duration  int     `json:"duration_minutes" validate:"required"`
}

// Handlers

// CreateCompany godoc
// @Summary Create a new company
// @Description Creates a new company and a default main branch
// @Tags companies
// @Accept json
// @Produce json
// @Param input body createCompanyRequest true "Company Input"
// @Success 201 {object} domain.Company
// @Failure 400 {object} erru.AppError
// @Router /companies [post]
func (h *CompanyHandler) CreateCompany(c echo.Context) error {
	var req createCompanyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	// Mock OwnerID from JWT (middleware should set this)
	// ownerID := c.Get("user_id").(uint)
	ownerID := uint(1) // Placeholder

	company, err := h.Usecase.CreateCompany(c.Request().Context(), req.Name, ownerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusCreated, company)
}

// AddBranch godoc
// @Summary Add a branch to a company
// @Tags companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body addBranchRequest true "Branch Input"
// @Success 201 {object} domain.Branch
// @Router /companies/{id}/branches [post]
func (h *CompanyHandler) AddBranch(c echo.Context) error {
	companyID, _ := strconv.Atoi(c.Param("id"))
	var req addBranchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	branch, err := h.Usecase.AddBranch(c.Request().Context(), uint(companyID), req.Name, req.Address, req.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusCreated, branch)
}

// AddCategory godoc
// @Summary Add a category
// @Tags companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body addCategoryRequest true "Category Input"
// @Success 201 {object} domain.Category
// @Router /companies/{id}/categories [post]
func (h *CompanyHandler) AddCategory(c echo.Context) error {
	companyID, _ := strconv.Atoi(c.Param("id"))
	var req addCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	cat, err := h.Usecase.AddCategory(c.Request().Context(), uint(companyID), req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusCreated, cat)
}

// AddService godoc
// @Summary Add a service to catalog
// @Tags companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body addServiceRequest true "Service Input"
// @Success 201 {object} domain.Service
// @Router /companies/{id}/services [post]
func (h *CompanyHandler) AddService(c echo.Context) error {
	companyID, _ := strconv.Atoi(c.Param("id"))
	var req addServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	svc, err := h.Usecase.AddService(c.Request().Context(), uint(companyID), req.CategoryID, req.Name, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusCreated, svc)
}

// AddEmployee godoc
// @Summary Add an employee to a branch
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Branch ID"
// @Param input body addEmployeeRequest true "Employee Input"
// @Success 201 {object} domain.Employee
// @Router /branches/{id}/employees [post]
func (h *CompanyHandler) AddEmployee(c echo.Context) error {
	branchID, _ := strconv.Atoi(c.Param("id"))
	var req addEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	emp, err := h.Usecase.AddEmployee(c.Request().Context(), uint(branchID), req.Name, req.Position)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusCreated, emp)
}

// AssignService godoc
// @Summary Assign service to employee with price
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param input body assignServiceRequest true "Assignment Input"
// @Success 200 {string} string "Assigned"
// @Router /employees/{id}/services [post]
func (h *CompanyHandler) AssignService(c echo.Context) error {
	employeeID, _ := strconv.Atoi(c.Param("id"))
	var req assignServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	err := h.Usecase.AssignService(c.Request().Context(), uint(employeeID), req.ServiceID, req.Price, req.Duration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "assigned"})
}

// GetEmployeeMenu godoc
// @Summary Get services performed by employee
// @Tags employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {array} domain.EmployeeService
// @Router /employees/{id}/services [get]
func (h *CompanyHandler) GetEmployeeMenu(c echo.Context) error {
	employeeID, _ := strconv.Atoi(c.Param("id"))

	menu, err := h.Usecase.GetEmployeeMenu(c.Request().Context(), uint(employeeID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, menu)
}
