package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/vipos89/timehub/pkg/erru"
	"github.com/vipos89/timehub/pkg/logger"
)

func PanicRecovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = erru.New(http.StatusInternalServerError, "Unknown panic")
				}
				stack := string(debug.Stack())
				logger.Error("Panic recovered", "error", err, "stack", stack)

				c.JSON(http.StatusInternalServerError, erru.ErrInternalServerError)
			}
		}()
		return next(c)
	}
}

func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Info("Request started",
			"method", c.Request().Method,
			"path", c.Request().URL.Path,
			"ip", c.RealIP(),
		)

		err := next(c)

		if err != nil {
			c.Error(err)
		}

		logger.Info("Request finished",
			"status", c.Response().Status,
		)

		return nil
	}
}
