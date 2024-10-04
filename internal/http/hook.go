package http

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
	"github.com/labstack/echo/v4"
)

func BeforeStart(e *echo.Echo, v core.Vault) {
	err := v.LoadCacheArchivedUserView.Invoke()

	if err != nil {
		logging.LogError("failed loading cache %v", err)
	}

	err = v.CacheArchivedUserView.ScheduleAll(func(username string) (model.ArchivedUserView, error) {
		return v.GetCachedArchivedUserView.Invoke(username)
	})

	if err != nil {
		logging.LogError("failed scheduling archive cache %v", err)
	}

	logging.LogInfo("Configured Routes:")
	routes := e.Routes()
	for _, r := range routes {
		logging.LogInfo("%s %s", r.Method, r.Path)
	}
}
