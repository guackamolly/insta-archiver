package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
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

func NewLoadCacheArchivedUserView(
	repository cache.CacheRepository[string, model.ArchivedUserView],
) LoadCacheArchivedUserView {
	return LoadCacheArchivedUserView{
		repository: repository,
	}
}

func NewCacheArchivedUserView(
	repository cache.CacheRepository[string, model.ArchivedUserView],
) CacheArchivedUserView {
	return CacheArchivedUserView{
		repository: repository,
	}
}

func NewGetCachedArchivedUserView(
	repository cache.CacheRepository[string, model.ArchivedUserView],
) GetCachedArchivedUserView {
	return GetCachedArchivedUserView{
		repository: repository,
	}
}
