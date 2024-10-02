package core

import "github.com/guackamolly/insta-archiver/internal/domain"

// Main application DI container.
type Vault struct {
	PurifyUsername            domain.PurifyUsername
	LoadCacheArchivedUserView domain.LoadCacheArchivedUserView
	CacheArchivedUserView     domain.CacheArchivedUserView
	GetCachedArchivedUserView domain.GetCachedArchivedUserView
	GetUserProfile            domain.GetUserProfile
	DownloadUserProfile       domain.DownloadUserProfile
}
