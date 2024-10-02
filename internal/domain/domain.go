package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

func NewGetUserProfile(
	repo user.UserRepository,
) GetUserProfile {
	return GetUserProfile{
		userRepository: repo,
	}
}

func NewDownloadUserProfile(
	client http.HttpClient,
	storage *storage.FileSystemStorage,
	physicalContentDir,
	virtualContentDir string,
) DownloadUserProfile {
	return DownloadUserProfile{
		client:             client,
		physicalContentDir: physicalContentDir,
		virtualContentDir:  virtualContentDir,
		contentStorage:     storage,
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
