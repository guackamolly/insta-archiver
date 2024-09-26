package http

import (
	"fmt"

	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/model"
	"github.com/labstack/echo/v4"
)

func onArchiveUserStories(
	ectx echo.Context,
	vault core.Vault,
	username string,
) model.ArchivedUserView {
	onError := ectx.Error

	pun := core.Invoke(username, vault.PurifyUsername, onError)
	view := core.Invoke(pun, vault.GetCachedArchivedUserView, onError)

	// If cached, return immediately
	if view != nil {
		fmt.Printf("user %s is cached, returning cached view\n", pun)
		return *view
	}

	cs := core.Invoke(pun, vault.GetLatestStories, onError)
	fs := core.Invoke(cs, vault.DownloadUserStories, onError)
	cs = core.Invoke(fs, vault.ArchiveUserStories, onError)
	cs = core.Invoke(cs, vault.PurifyCloudStories, onError)

	v := model.NewArchivedUserView(pun, "Something along these lines", cs)

	core.Invoke(
		v,
		vault.CacheArchivedUserView,
		func(err error) {
			// if cache fails, we can still rely on memory value
			ectx.Logger().Error(err)
		},
	)

	return v
}
