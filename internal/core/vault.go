package core

import "github.com/guackamolly/insta-archiver/internal/domain"

// Main application DI container.
type Vault struct {
	DownloadUserStories       domain.DownloadUserStories
	GetLatestStories          domain.GetLatestStories
	ArchiveUserStories        domain.ArchiveUserStories
	PurifyCloudStories        domain.PurifyCloudStories
	PurifyUsername            domain.PurifyUsername
	LoadCacheArchivedUserView domain.LoadCacheArchivedUserView
	CacheArchivedUserView     domain.CacheArchivedUserView
	GetCachedArchivedUserView domain.GetCachedArchivedUserView
}
