package domain

import "github.com/guackamolly/insta-archiver/internal/model"

const (
	ArchiveFailed model.ErrorReason = iota + 1
	LoadCacheFailed
	UpdateCacheFailed
	LookupCacheFailed
	DownloadThumbnailFailed
	DownloadMediaFailed
	DownloadAvatarFailed
	FetchBioFailed
	StoreStoriesFailed
	StoreAvatarFailed
	FetchProfileFailed
	FetchStoriesFailed
	ConvertStoriesFailed
	ValidateUsernameFailed
	ScheduleArchiveFailed
)
