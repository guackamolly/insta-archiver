package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(loggingMiddleware())
}

func loggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			req := ectx.Request()

			println(fmt.Sprintf("Host: %s | Method: %s | Path: %s | Client IP: %s", req.Host, req.Method, req.URL.RequestURI(), ectx.RealIP()))
			return next(ectx)
		}
	}
}
