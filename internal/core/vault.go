package core

import "github.com/guackamolly/insta-archiver/internal/domain"

// Main application DI container.
type Vault struct {
	DownloadUserAvatar        domain.DownloadUserAvatar
	ArchiveUserAvatar         domain.ArchiveUserAvatar
	DownloadUserStories       domain.DownloadUserStories
	FilterStoriesForDownload  domain.FilterStoriesForDownload
	GetArchivedStories        domain.GetArchivedStories
	GetLatestStories          domain.GetLatestStories
	GetUserBio                domain.GetUserBio
	ArchiveUserStories        domain.ArchiveUserStories
	PurifyStaticUrl           domain.PurifyStaticUrl
	PurifyUsername            domain.PurifyUsername
	LoadCacheArchivedUserView domain.LoadCacheArchivedUserView
	CacheArchivedUserView     domain.CacheArchivedUserView
	GetCachedArchivedUserView domain.GetCachedArchivedUserView
}
