package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/model"
)

func NewGetUserProfile(
	archiveRepo archive.ArchiveRepository,
	userRepo user.UserRepository,
) GetUserProfile {
	return GetUserProfile{
		archiveRepository: archiveRepo,
		userRepository:    userRepo,
	}
}

func NewDownloadUserProfile(
	archiveRepo archive.ArchiveRepository,
) DownloadUserProfile {
	return DownloadUserProfile{
		archiveRepo: archiveRepo,
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
	archiveRepo archive.ArchiveRepository,
	cacheRepo cache.CacheRepository[string, model.ArchivedUserView],
) CacheArchivedUserView {
	return CacheArchivedUserView{
		archiveRepo: archiveRepo,
		cacheRepo:   cacheRepo,
	}
}

func NewGetCachedArchivedUserView(
	repository cache.CacheRepository[string, model.ArchivedUserView],
) GetCachedArchivedUserView {
	return GetCachedArchivedUserView{
		repository: repository,
	}
}
