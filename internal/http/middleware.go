package http

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/labstack/echo/v4"
)

const vaultKey = "vault"

func RegisterMiddlewares(e *echo.Echo, vault core.Vault) {
	e.Use(loggingMiddleware())
	e.Use(vaultMiddleware(vault))
}

func vaultMiddleware(vault core.Vault) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			ectx.Set(vaultKey, vault)

			return next(ectx)
		}
	}
}

func loggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			req := ectx.Request()

			logging.LogInfo("Host: %s | Method: %s | Path: %s | Client IP: %s", req.Host, req.Method, req.URL.RequestURI(), ectx.RealIP())
			return next(ectx)
		}
	}
}

func withVault(ectx echo.Context, with func(core.Vault) error) error {
	v, ok := ectx.Get(vaultKey).(core.Vault)

	if ok {
		return with(v)
	}

	return echo.ErrFailedDependency
}
