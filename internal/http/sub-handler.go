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

	// 1. Validate username
	pun, err := domain.Invoke(username, vault.PurifyUsername, nil)
	// 2. Check if view is cached
	view, cerr := domain.Invoke(pun, vault.GetCachedArchivedUserView, err)

	// If cached, return immediately
	if cerr == nil && view != nil {
		fmt.Printf("user %s is cached, returning cached view\n", pun)
		return *view, nil
	}

	// 3. Get user bio
	bio, berr := domain.Invoke(pun, vault.GetUserBio, err)

	// if bio fetch failed, use the default one
	if err == nil && berr != nil {
		bio = model.DefaultBio()
	}

	if err == nil && berr == nil {
		df, dferr := domain.Invoke(bio, vault.DownloadUserAvatar, berr)

		if dferr == nil {
			aa, aaerr := vault.ArchiveUserAvatar.Invoke(pun, df)
			if aaerr == nil {
				url, _ := vault.PurifyStaticUrl.Invoke(aa.Path)
				bio = model.NewBio(bio.Avatar, bio.Description, url)
			}
		}
	}

	// 4. Execute archiving calls (Fetch + Download + Archive + Cache)
	cs, err := domain.Invoke(pun, vault.GetLatestStories, err)
	fcs, aerr := domain.Invoke(cs, vault.FilterStoriesForDownload, err)

	if aerr == nil {
		cs = fcs
	}

	fs, err := domain.Invoke(cs, vault.DownloadUserStories, err)
	_, err = domain.Invoke(fs, vault.ArchiveUserStories, err)
	cs, err = domain.Invoke(username, vault.GetArchivedStories, err)

	if err != nil {
		return model.ArchivedUserView{}, err
	}

	cs, err = vault.PurifyStaticUrl.InvokeStories(cs)
	v, cerr := domain.Invoke(
		model.NewArchivedUserView(
			pun,
			bio,
			cs,
		),
		vault.CacheArchivedUserView,
		err,
	)

	// if cache fails, we can still rely on memory value
	if err == nil && cerr != nil {
		ectx.Logger().Error(err)
	}

	return v, err
}
