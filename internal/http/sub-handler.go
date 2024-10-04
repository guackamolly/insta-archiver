package http

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
)

func ArchiveUser(
	vault core.Vault,
	username string,
) (model.ArchivedUserView, error) {
	var view model.ArchivedUserView

	// 1. Validate username
	pun, err := vault.PurifyUsername.Invoke(username)

	if err != nil {
		return view, err
	}

	// 2. Check if view is cached
	view, err = vault.GetCachedArchivedUserView.Invoke(pun)

	// If cached, return immediately
	if err == nil {
		logging.LogInfo("user %s is cached, returning cached view", pun)
		return view, nil
	}

	view, err = getAndCacheUserProfile(pun, vault)

	if err == nil && !view.IsPrivate {
		vault.CacheArchivedUserView.Schedule(view.Username, func() (model.ArchivedUserView, error) {
			return getAndCacheUserProfile(view.Username, vault)
		})
	}

	return view, err
}

func getAndCacheUserProfile(
	username string,
	vault core.Vault,
) (model.ArchivedUserView, error) {
	// 1. Get user profile
	profile, err := vault.GetUserProfile.Invoke(username)

	if err != nil {
		return model.ArchivedUserView{}, err
	}

	// 2. Download user profile
	profile, err = vault.DownloadUserProfile.Invoke(profile)

	if err != nil {
		return model.ArchivedUserView{}, err
	}

	// 3. Cache view
	view := model.NewArchivedUserView(profile)
	view, err = vault.CacheArchivedUserView.Invoke(view)

	return view, err
}
