package http

import (
	"net/http"

	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.Any(rootRoute, anyRouteHandler)
	e.GET(archiveRoute, archiveRouteHandler)
	echo.NotFoundHandler = useNotFoundHandler()
}

func archiveRouteHandler(ectx echo.Context) error {
	return withVault(ectx, func(v core.Vault) error {
		un := ectx.QueryParam(archiveQueryParam)
		resp := onArchiveUserStories(ectx, v, un)

		return ectx.Render(http.StatusOK, "index.html", resp)
	})
}

func anyRouteHandler(ectx echo.Context) error {
	ap := ectx.QueryParam(archiveQueryParam)

	if len(ap) == 0 {
		return ectx.File(root)
	} else {
		return archiveRouteHandler(ectx)
	}
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.File(errors[404])
	}
}
