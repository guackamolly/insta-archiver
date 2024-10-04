package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type DownloadUserProfile struct {
	archiveRepo archive.ArchiveRepository
}

func (u DownloadUserProfile) Invoke(profile model.Profile) (model.Profile, error) {
	return u.archiveRepo.Archive(profile)
}
