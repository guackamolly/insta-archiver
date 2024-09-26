package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type ArchiveUserStories struct {
	repository archive.ArchiveRepository
}

func (u ArchiveUserStories) Invoke(stories []model.FileStory) ([]model.CloudStory, error) {
	return u.repository.Archive(stories...)
}
