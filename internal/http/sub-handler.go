package http

import (
	"fmt"

	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/domain"
	"github.com/guackamolly/insta-archiver/internal/model"
	"github.com/labstack/echo/v4"
)

func onArchiveUserStories(
	ectx echo.Context,
	vault core.Vault,
	username string,
) (model.ArchivedUserView, error) {

	// 1. Validate username + check if view is cached
	pun, err := domain.Invoke(username, vault.PurifyUsername, nil)
	view, cerr := domain.Invoke(pun, vault.GetCachedArchivedUserView, err)

	// If cached, return immediately
	if cerr == nil && view != nil {
		fmt.Printf("user %s is cached, returning cached view\n", pun)
		return *view, nil
	}

	// 2. Execute archiving callsd (Fetch + Download + Archive + Cache)
	cs, err := domain.Invoke(pun, vault.GetLatestStories, err)
	fs, err := domain.Invoke(cs, vault.DownloadUserStories, err)
	cs, err = domain.Invoke(fs, vault.ArchiveUserStories, err)
	cs, err = domain.Invoke(cs, vault.PurifyCloudStories, err)
	v, cerr := domain.Invoke(model.NewArchivedUserView(pun, "Something along these lines", cs), vault.CacheArchivedUserView, err)

	// if cache fails, we can still rely on memory value
	if err == nil && cerr != nil {
		ectx.Logger().Error(err)
	}

	return v, err
}
