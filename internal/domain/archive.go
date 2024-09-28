package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type ArchiveUserStories struct {
	repository archive.ArchiveRepository
}

func (u ArchiveUserStories) Invoke(stories []model.FileStory) ([]model.CloudStory, error) {
	return WrapResult(stories, u.repository.Archive, ArchiveFailed)
}

type ArchiveUserAvatar struct {
	storage *storage.FileSystemStorage
}

func (u ArchiveUserAvatar) Invoke(username string, avatar storage.File) (storage.File, error) {
	af, err := u.storage.Store(username, []storage.File{avatar})

	if err == nil {
		return af[0], nil
	}

	var f storage.File

	return f, model.Wrap(err, ArchiveFailed)
}
