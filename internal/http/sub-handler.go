package http

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/model"
	"github.com/labstack/echo/v4"
)

func onArchiveUserStories(
	ectx echo.Context,
	vault core.Vault,
	username string,
) model.ArchivedUserView {
	onError := func(err error) error { return onError(ectx, err) }

	pun := core.Invoke(username, vault.PurifyUsername, onError)
	view := core.Invoke(pun, vault.GetCachedArchivedUserView, onError)

	// If cached, return immediately
	if view != nil {
		return *view
	}

	cs := core.Invoke(pun, vault.GetLatestStories, onError)
	fs := core.Invoke(cs, vault.DownloadUserStories, onError)
	cs = core.Invoke(fs, vault.ArchiveUserStories, onError)
	cs = core.Invoke(cs, vault.PurifyCloudStories, onError)

	// cache view on return
	return core.Invoke(
		model.NewArchivedUserView(pun, "Something along these lines", cs),
		vault.CacheArchivedUserView,
		onError,
	)
}

func onError(
	ectx echo.Context,
	err error,
) error {
	return ectx.String(400, err.Error())
}
