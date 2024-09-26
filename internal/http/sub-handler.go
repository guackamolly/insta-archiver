package http

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/labstack/echo/v4"
)

func onArchiveUserStories(
	ectx echo.Context,
	vault core.Vault,
	username string,
) ArchiveUserStoriesResponse {
	onError := func(err error) error { return onError(ectx, err) }

	pun := core.Invoke(username, vault.PurifyUsername, onError)
	cs := core.Invoke(pun, vault.GetLatestStories, onError)
	fs := core.Invoke(cs, vault.DownloadUserStories, onError)
	cs = core.Invoke(fs, vault.ArchiveUserStories, onError)
	cs = core.Invoke(cs, vault.PurifyCloudStories, onError)

	return ArchiveUserStoriesResponse{
		Username:    pun,
		LastStories: cs,
	}
}

func onError(
	ectx echo.Context,
	err error,
) error {
	return ectx.String(400, err.Error())
}
