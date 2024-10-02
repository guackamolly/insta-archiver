package http

import (
	"time"

	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
	"github.com/labstack/echo/v4"
)

func BeforeStart(e *echo.Echo, v core.Vault) {
	vs, err := v.LoadCacheArchivedUserView.Invoke()

	if err != nil {
		logging.LogError("failed loading cache %v", err)

		return
	}

	for _, ce := range vs {
		view := ce.Value
		v.CacheArchivedUserView.Schedule(view.Username, time.Until(ce.NextHit), func() (model.ArchivedUserView, error) {
			return getAndCacheUserProfile(view.Username, v)
		})
	}
}
