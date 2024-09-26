package http

import (
	"net/http"

	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.Any(rootRoute, anyRouteHandler)
	e.GET(archiveRoute, archiveRouteHandler)
	e.HTTPErrorHandler = httpErrorHandler()
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

func httpErrorHandler() func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		he, ok := err.(*echo.HTTPError)

		// If cast fails, serve fallback
		if !ok {
			c.File(fallback)
			return
		}

		// if error page available, serve it
		if f, ok := errors[he.Code]; !ok {
			c.File(f)
			return
		}

		// If all checks fail, serve fallback
		c.File(fallback)
	}
}
