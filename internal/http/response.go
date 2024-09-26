package http

import "github.com/guackamolly/insta-archiver/internal/model"

type ArchiveUserStoriesResponse struct {
	Username        string
	Description     string
	LastStories     []model.CloudStory
	ArchivedStories []model.CloudStory
}
