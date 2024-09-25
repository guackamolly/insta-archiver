package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type GetLatestStories struct {
	repository user.UserRepository
}

func (u GetLatestStories) Invoke(username string) ([]model.CloudStory, error) {
	return u.repository.Stories(username)
}
