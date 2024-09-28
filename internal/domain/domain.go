package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

func NewFilterStoriesForDownload(
	repo archive.ArchiveRepository,
) FilterStoriesForDownload {
	return FilterStoriesForDownload{
		repository: repo,
	}
}

func NewGetArchivedStories(
	repo archive.ArchiveRepository,
) GetArchivedStories {
	return GetArchivedStories{
		repository: repo,
	}
}

func NewGetUserBio(
	cacheRepo cache.MemoryCacheRepository[model.Bio],
	userRepo user.UserRepository,
) GetUserBio {
	return GetUserBio{
		userRepository:  userRepo,
		cacheRepository: cacheRepo,
	}
}

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

func NewDownloadUserAvatar(
	client http.HttpClient,
) DownloadUserAvatar {
	return DownloadUserAvatar{
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

func NewArchiveUserAvatar(
	storage *storage.FileSystemStorage,
) ArchiveUserAvatar {
	return ArchiveUserAvatar{
		storage: storage,
	}
}

func NewPurifyStaticUrl(
	physicalContentDir,
	virtualContentDir string,
) PurifyStaticUrl {
	return PurifyStaticUrl{
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
