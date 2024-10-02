package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type GetUserProfile struct {
	userRepository user.UserRepository
}

func (u GetUserProfile) Invoke(username string) (model.Profile, error) {
	logging.LogInfo("getting bio for %s...", username)
	bio, err := u.userRepository.Bio(username)
	if err != nil {
		return model.Profile{}, model.Wrap(err, FetchBioFailed)
	}

	logging.LogInfo("getting stories for %s...", username)
	stories, err := u.userRepository.Stories(username)
	if err != nil {
		return model.Profile{}, model.Wrap(err, FetchStoriesFailed)
	}

	return model.NewProfile(bio, stories), nil
}
