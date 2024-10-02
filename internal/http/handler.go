package http

import (
	"net/http"

	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.Any(rootRoute, anyRouteHandler)
	e.GET(archiveRoute, archiveRouteHandler)
	e.HTTPErrorHandler = httpErrorHandler()
}

func archiveRouteHandler(ectx echo.Context) error {
	un := ectx.QueryParam(archiveQueryParam)
	if un == "" {
		un = ectx.Param("id")
	}

	return withVault(ectx, func(v core.Vault) error {
		resp, err := ArchiveUser(v, un)

		if err != nil {
			return err
		}

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
		// make sure to not process any false positives
		if err == nil {
			return
		}

		logging.LogError("handling error... %v", err)

		me, ok := err.(*model.Error)

		if ok {
			onCustomError(c, me)
			return
		}

		he, ok := err.(*echo.HTTPError)

		// If all cast fail, serve fallback
		if !ok {
			c.File(fallback)
			return
		}

		// if error page available, serve it
		if f, eok := errors[he.Code]; eok {
			c.File(f)
			return
		}

		// if no match, resort to fallback
		c.File(fallback)
	}
}
