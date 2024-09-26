package archive

import (
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Implements [ArchiveRepository] using [FileSystemStorage] as the data source.
type FileSystemArchiveRepository struct {
	storage *storage.FileSystemStorage
}

func (r FileSystemArchiveRepository) Archive(stories ...model.FileStory) ([]model.CloudStory, error) {
	cs := make([]model.CloudStory, len(stories))

	for i, v := range stories {
		fs, err := r.storage.Store(v.Username, []storage.File{v.Thumbnail, v.Media})

		if err != nil {
			return cs, err
		}

		cs[i] = model.CloudStory{
			Id:        v.Id,
			Username:  v.Username,
			Thumbnail: fs[0].Path,
			Media:     fs[1].Path,
		}
	}

	return cs, nil
}