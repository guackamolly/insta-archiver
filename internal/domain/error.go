package domain

import "github.com/guackamolly/insta-archiver/internal/model"

const (
	LoadCacheFailed model.ErrorReason = iota + 1
	UpdateCacheFailed
	LookupCacheFailed
	FetchBioFailed
	FetchStoriesFailed
	ValidateUsernameFailed
	ScheduleArchiveFailed
)
