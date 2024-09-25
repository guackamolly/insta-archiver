package http

import (
	"net/http"

	_http "github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/user"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.Any(rootRoute, anyRouteHandler)
	e.GET(archiveRoute, archiveRouteHandler)
	echo.NotFoundHandler = useNotFoundHandler()
}

func archiveRouteHandler(ectx echo.Context) error {
	r := user.ViewIGStoryUserRepository{
		Client: _http.Native(),
	}

	ap := ectx.QueryParam(archiveQueryParam)

	res, _ := r.Stories(ap)

	return ectx.Render(http.StatusOK, "index.html", map[string]any{
		"username":             ap,
		"description":          "I love Messi. SIUUUUUUU",
		"archivedStoriesCount": len(res),
		"lastStoriesCount":     len(res),
		"stories":              res,
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
