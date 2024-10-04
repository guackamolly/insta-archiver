package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type GetUserProfile struct {
	userRepository    user.UserRepository
	archiveRepository archive.ArchiveRepository
}

func (u GetUserProfile) Invoke(username string) (model.Profile, error) {
	logging.LogInfo("getting bio for %s...", username)
	bio, err := u.userRepository.Bio(username)
	if err != nil {
		return model.Profile{}, model.Wrap(err, FetchBioFailed)
	}

	logging.LogInfo("getting stories for %s...", username)
	lastStories, lerr := u.userRepository.Stories(username)
	if lerr != nil {
		logging.LogWarning("failed to fetch last stories... %v", lerr)
	}

	archivedStories, err := u.archiveRepository.Stories(username)
	if err != nil && lerr != nil {
		logging.LogWarning("failed to fetch archived stories... %v", err)
		return model.Profile{}, model.Wrap(err, FetchStoriesFailed)
	}

	stories := archivedStories

	for _, ls := range lastStories {
		dup := false
		for _, as := range archivedStories {
			if ls.Id == as.Id {
				dup = true
				break
			}
		}

		if !dup {
			stories = append(stories, ls)
		}
	}

	return model.NewProfile(bio, stories), nil
}
