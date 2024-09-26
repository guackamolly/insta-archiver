package http

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/labstack/echo/v4"
)

func BeforeStart(e *echo.Echo, v core.Vault) {
	err := v.LoadCacheArchivedUserView.Invoke()

	if err != nil {
		e.Logger.Error(err)
	}
}
