package http

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
)

func onArchiveUserStories(
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

	// 3. Get user profile
	profile, err := vault.GetUserProfile.Invoke(pun)

	if err != nil {
		return view, err
	}

	// 4. Download user profile
	profile, err = vault.DownloadUserProfile.Invoke(profile)

	if err != nil {
		return view, err
	}

	//5. Cache view
	view = model.NewArchivedUserView(profile)
	view, err = vault.CacheArchivedUserView.Invoke(view)

	return view, err
}
