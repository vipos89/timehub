package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vipos89/timehub/pkg/erru"
	"github.com/vipos89/timehub/services/company-service/internal/domain"
)

type CompanyHandler struct {
	Usecase domain.CompanyUsecase
}

func NewCompanyHandler(e *echo.Echo, us domain.CompanyUsecase) {
	handler := &CompanyHandler{
		Usecase: us,
	}

	// Company Routes
	e.POST("/companies", handler.CreateCompany)
	e.GET("/companies", handler.GetCompanies)
	e.GET("/companies/:id", handler.GetCompanyByID)
	e.POST("/companies/:id/branches", handler.AddBranch)
	e.GET("/companies/:id/branches", handler.GetBranches)

	// Branch-Specific Service/Category Routes
	e.POST("/branches/:id/categories", handler.AddCategory)
	e.GET("/branches/:id/categories", handler.GetCategories)
	e.POST("/branches/:id/services", handler.AddService)
	e.GET("/branches/:id/services", handler.GetServices)
	e.PUT("/services/:id", handler.UpdateService)

	// Employee Routes
	e.GET("/employees", handler.GetEmployees) // Query param company_id
	e.POST("/employees", handler.AddEmployee)
	e.POST("/branches/:id/employees", handler.AddEmployee) // Keep for backward compatibility
	e.POST("/employees/:id/services", handler.AssignService)
	e.DELETE("/employees/:id/services/:serviceId", handler.RemoveService)
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
	CompanyID       uint    `json:"company_id"` // Optional or from path
	CategoryID      *uint   `json:"category_id"`
	Name            string  `json:"name" validate:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
}

type updateServiceRequest struct {
	CategoryID      *uint   `json:"category_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
}

type addEmployeeRequest struct {
	BranchID uint   `json:"branch_id"`
	Name     string `json:"name" validate:"required"`
	Position string `json:"position"`
	Email    string `json:"email"` // For Auth Service registration
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

// GetCompanies godoc
// @Summary Get all companies owned by user
// @Tags companies
// @Produce json
// @Success 200 {array} domain.Company
// @Router /companies [get]
func (h *CompanyHandler) GetCompanies(c echo.Context) error {
	// ownerID := c.Get("user_id").(uint)
	ownerID := uint(1) // Placeholder

	companies, err := h.Usecase.GetMyCompanies(c.Request().Context(), ownerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, companies)
}

// GetCompanyByID godoc
// @Summary Get company details by ID
// @Tags companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} domain.Company
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetCompanyByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	company, err := h.Usecase.GetCompanyByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	if company == nil {
		return c.JSON(http.StatusNotFound, erru.ErrNotFound)
	}
	return c.JSON(http.StatusOK, company)
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

// GetBranches godoc
// @Summary Get branches of a company
// @Tags companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {array} domain.Branch
// @Router /companies/{id}/branches [get]
func (h *CompanyHandler) GetBranches(c echo.Context) error {
	companyID, _ := strconv.Atoi(c.Param("id"))
	branches, err := h.Usecase.GetCompanyBranches(c.Request().Context(), uint(companyID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, branches)
}

// AddCategory godoc
// @Summary Add a category to a branch
// @Tags companies
// @Accept json
// @Produce json
// @Param id path int true "Branch ID"
// @Param input body addCategoryRequest true "Category Input"
// @Success 201 {object} domain.Category
// @Router /branches/{id}/categories [post]
func (h *CompanyHandler) AddCategory(c echo.Context) error {
	branchID, _ := strconv.Atoi(c.Param("id"))
	var req addCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	// In a real app, we'd get CompanyID from the branch itself or JWT
	// For now, let's assume companyID 1 or fetch it
	companyID := uint(1)

	cat, err := h.Usecase.AddCategory(c.Request().Context(), companyID, uint(branchID), req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusCreated, cat)
}

// GetCategories godoc
// @Summary Get categories of a branch
// @Tags companies
// @Produce json
// @Param id path int true "Branch ID"
// @Success 200 {array} domain.Category
// @Router /branches/{id}/categories [get]
func (h *CompanyHandler) GetCategories(c echo.Context) error {
	branchID, _ := strconv.Atoi(c.Param("id"))
	cats, err := h.Usecase.GetBranchCategories(c.Request().Context(), uint(branchID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, cats)
}

// AddService godoc
// @Summary Add a service to a branch catalog
// @Tags companies
// @Accept json
// @Produce json
// @Param id path int true "Branch ID"
// @Param input body addServiceRequest true "Service Input"
// @Success 201 {object} domain.Service
// @Router /branches/{id}/services [post]
func (h *CompanyHandler) AddService(c echo.Context) error {
	branchID, _ := strconv.Atoi(c.Param("id"))
	var req addServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	companyID := req.CompanyID
	if companyID == 0 {
		companyID = uint(1)
	}

	svc, err := h.Usecase.AddService(c.Request().Context(), companyID, uint(branchID), req.CategoryID, req.Name, req.Description, req.Price, req.DurationMinutes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusCreated, svc)
}

// UpdateService godoc
// @Summary Update service details
// @Tags companies
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param input body updateServiceRequest true "Service Update Input"
// @Success 200 {object} domain.Service
// @Router /services/{id} [put]
func (h *CompanyHandler) UpdateService(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req updateServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	// For simplicity, we just pass the object to usecase
	// In a real app, we'd fetch it first to ensure it exists and belongs to the company
	svc := &domain.Service{
		ID:              uint(id),
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
	}

	err := h.Usecase.UpdateService(c.Request().Context(), svc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, svc)
}

func (h *CompanyHandler) RemoveService(c echo.Context) error {
	employeeID, _ := strconv.Atoi(c.Param("id"))
	serviceID, _ := strconv.Atoi(c.Param("serviceId"))

	err := h.Usecase.RemoveService(c.Request().Context(), uint(employeeID), uint(serviceID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}

// GetServices godoc
// @Summary Get all services of a branch
// @Tags companies
// @Produce json
// @Param id path int true "Branch ID"
// @Success 200 {array} domain.Service
// @Router /branches/{id}/services [get]
func (h *CompanyHandler) GetServices(c echo.Context) error {
	branchID, _ := strconv.Atoi(c.Param("id"))
	svcs, err := h.Usecase.GetBranchServices(c.Request().Context(), uint(branchID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, svcs)
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
	var req addEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	// Support both path param and body param for branchID
	branchID := req.BranchID
	if branchID == 0 {
		id, _ := strconv.Atoi(c.Param("id"))
		branchID = uint(id)
	}

	if branchID == 0 {
		return c.JSON(http.StatusBadRequest, erru.New(http.StatusBadRequest, "branch_id is required"))
	}

	emp, err := h.Usecase.AddEmployee(c.Request().Context(), branchID, req.Name, req.Position, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusCreated, emp)
}

// GetEmployees godoc
// @Summary Get employees of a company
// @Tags employees
// @Produce json
// @Param company_id query int true "Company ID"
// @Success 200 {array} domain.Employee
// @Router /employees [get]
func (h *CompanyHandler) GetEmployees(c echo.Context) error {
	companyID, _ := strconv.Atoi(c.QueryParam("company_id"))
	emps, err := h.Usecase.GetCompanyEmployees(c.Request().Context(), uint(companyID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, emps)
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
