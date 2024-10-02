package domain

import (
	"fmt"

	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type GetUserProfile struct {
	userRepository user.UserRepository
}

func (u GetUserProfile) Invoke(username string) (model.Profile, error) {
	fmt.Printf("getting bio for %s...\n", username)
	bio, err := u.userRepository.Bio(username)
	if err != nil {
		return model.Profile{}, model.Wrap(err, FetchBioFailed)
	}

	fmt.Printf("getting stories for %s...\n", username)
	stories, err := u.userRepository.Stories(username)
	if err != nil {
		return model.Profile{}, model.Wrap(err, FetchStoriesFailed)
	}

	return model.NewProfile(bio, stories), nil
}
