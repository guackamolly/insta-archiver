package domain

import (
	"strings"

	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type GetArchivedStories struct {
	repository archive.ArchiveRepository
}

type GetLatestStories struct {
	repository user.UserRepository
}

type FilterStoriesForDownload struct {
	repository archive.ArchiveRepository
}

type PurifyCloudStories struct {
	physicalContentDir string
	virtualContentDir  string
}

func (u FilterStoriesForDownload) Invoke(stories []model.CloudStory) ([]model.CloudStory, error) {
	if len(stories) == 0 {
		return stories, nil
	}

	username := stories[0].Username
	as, err := u.repository.All(username)

	if err != nil {
		return stories, model.Wrap(err, FetchStoriesFailed)
	}

	return model.Filter(stories, func(s model.CloudStory) bool {
		return model.Find(as, func(ss model.CloudStory) bool {
			return model.SameStory(s, ss)
		}) == nil
	}), nil
}

func (u GetArchivedStories) Invoke(username string) ([]model.CloudStory, error) {
	return WrapResult(username, u.repository.All, FetchStoriesFailed)
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
