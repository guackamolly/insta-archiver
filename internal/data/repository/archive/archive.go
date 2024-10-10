package archive

import (
	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type ArchiveRepository interface {
	All() ([]string, error)
	Stories(username string) ([]model.Story, error)
	Archive(profile model.Profile) (model.Profile, error)
}

func NewFileSystemArchiveRepository(
	storage *storage.FileSystemStorage,
	client http.HttpClient,
	virtualContentDir string,
) ArchiveRepository {
	return FileSystemArchiveRepository{
		storage:           storage,
		client:            client,
		virtualContentDir: virtualContentDir,
	}
}

func NewFileSystemCDNArchiveRepository(
	storage *storage.FileSystemStorage,
	client http.HttpClient,
	cdnUrl string,
) ArchiveRepository {
	return FileSystemCDNArchiveRepository{
		storage: storage,
		client:  client,
		cdnUrl:  cdnUrl,
	}
}
