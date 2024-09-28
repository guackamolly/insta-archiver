package archive

import (
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Data operations related to stories archiving
type ArchiveRepository interface {
	All(username string) ([]model.CloudStory, error)
	Archive(stories []model.FileStory) ([]model.CloudStory, error)
}

func NewFileSystemArchiveRepository(
	storage *storage.FileSystemStorage,
) ArchiveRepository {
	return FileSystemArchiveRepository{
		storage: storage,
	}
}
