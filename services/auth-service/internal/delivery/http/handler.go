package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vipos89/timehub/pkg/erru"
	"github.com/vipos89/timehub/services/auth-service/internal/domain"
	"github.com/vipos89/timehub/services/auth-service/internal/usecase"
)

type AuthHandler struct {
	AuthUsecase usecase.AuthUsecase
}

func NewAuthHandler(e *echo.Echo, us usecase.AuthUsecase) {
	handler := &AuthHandler{
		AuthUsecase: us,
	}

	e.POST("/auth/register", handler.Register)
	e.POST("/auth/login", handler.Login)
}

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required"` // owner, admin, master
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user (company owner or employee)
// @Tags auth
// @Accept json
// @Produce json
// @Param input body registerRequest true "Register Input"
// @Success 201 {object} domain.User
// @Failure 400 {object} erru.AppError
// @Failure 500 {object} erru.AppError
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	user, err := h.AuthUsecase.Register(c.Request().Context(), req.Email, req.Password, req.Role)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return c.JSON(http.StatusConflict, erru.New(http.StatusConflict, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, erru.New(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusCreated, user)
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body loginRequest true "Login Input"
// @Success 200 {object} loginResponse
// @Failure 400 {object} erru.AppError
// @Failure 401 {object} erru.AppError
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, erru.ErrBadRequest)
	}

	token, err := h.AuthUsecase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, erru.ErrUnauthorized)
	}

	return c.JSON(http.StatusOK, loginResponse{Token: token})
}
