package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

func NewGetLatestStories(
	repo user.UserRepository,
) GetLatestStories {
	return GetLatestStories{
		repository: repo,
	}
}

func NewDownloadUserStories(
	client http.HttpClient,
) DownloadUserStories {
	return DownloadUserStories{
		client: client,
	}
}

func NewArchiveUserStories(
	repo archive.ArchiveRepository,
) ArchiveUserStories {
	return ArchiveUserStories{
		repository: repo,
	}
}

func NewPurifyCloudStories(
	physicalContentDir,
	virtualContentDir string,
) PurifyCloudStories {
	return PurifyCloudStories{
		physicalContentDir: physicalContentDir,
		virtualContentDir:  virtualContentDir,
	}
}

func NewPurifyUsername() PurifyUsername {
	return PurifyUsername{}
}

func NewCacheArchivedUserView(
	storage *storage.MemoryStorage[string, model.ArchivedUserView],
) CacheArchivedUserView {
	return CacheArchivedUserView{
		storage: storage,
	}
}

func NewGetCachedArchivedUserView(
	storage *storage.MemoryStorage[string, model.ArchivedUserView],
) GetCachedArchivedUserView {
	return GetCachedArchivedUserView{
		storage: storage,
	}
}
