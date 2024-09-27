package domain

import (
	"strings"

	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type GetLatestStories struct {
	repository user.UserRepository
}

type PurifyCloudStories struct {
	physicalContentDir string
	virtualContentDir  string
}

func (u GetLatestStories) Invoke(username string) ([]model.CloudStory, error) {
	return WrapResult(username, u.repository.Stories, FetchStoriesFailed)
}

func (u PurifyCloudStories) Invoke(stories []model.CloudStory) ([]model.CloudStory, error) {
	cs := make([]model.CloudStory, len(stories))

	for i, s := range stories {
		s.Thumbnail = u.purifyUrl(s.Thumbnail)
		s.Media = u.purifyUrl(s.Media)
		cs[i] = s
	}

	return cs, nil
}

func (u PurifyCloudStories) purifyUrl(url string) string {
	if strings.HasPrefix(url, u.physicalContentDir) {
		return strings.Replace(url, u.physicalContentDir, u.virtualContentDir, 1)
	}

	return url
}
