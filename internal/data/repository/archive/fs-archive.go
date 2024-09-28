package archive

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Implements [ArchiveRepository] using [FileSystemStorage] as the data source.
type FileSystemArchiveRepository struct {
	storage *storage.FileSystemStorage
}

func (r FileSystemArchiveRepository) All(username string) ([]model.CloudStory, error) {
	fs, err := r.storage.Lookup(username)

	if err != nil {
		return nil, err
	}

	return model.MapFilter(
		fs,
		func(f storage.File) (model.CloudStory, error) {
			var s model.CloudStory

			bs, err := os.ReadFile(f.Path)

			if err != nil {
				return s, err
			}

			err = json.Unmarshal(bs, &s)
			return s, err
		},
		func(f storage.File) bool {
			return strings.HasSuffix(f.Path, ".dat")
		},
	), nil
}

func (r FileSystemArchiveRepository) Archive(stories []model.FileStory) ([]model.CloudStory, error) {
	cs := make([]model.CloudStory, len(stories))

	for i, v := range stories {
		fs, err := r.storage.Store(v.Username, []storage.File{v.Thumbnail, v.Media})

		if err != nil {
			return cs, err
		}

		cs[i] = model.NewStory(
			v.Id,
			v.Username,
			v.PublishedOn,
			v.IsVideo,
			fs[0].Path,
			fs[1].Path,
		)

		bs, err := json.Marshal(cs[i])

		if err != nil {
			return cs, err
		}

		_, err = r.storage.StoreRaw(fmt.Sprintf("%s/%s.dat", v.Username, v.Id), bs)

		if err != nil {
			return cs, err
		}
	}

	return cs, nil
}
