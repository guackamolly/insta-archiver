package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type DownloadUserStories struct {
	client http.HttpClient
}

func (u DownloadUserStories) Invoke(stories []model.CloudStory) ([]model.FileStory, error) {
	fs := make([]model.FileStory, len(stories))

	for i, v := range stories {
		t, terr := u.downloadAndStat(v.Thumbnail)
		m, merr := u.downloadAndStat(v.Media)

		if terr != nil {
			return nil, terr
		}

		if merr != nil {
			return nil, merr
		}

		fs[i] = model.NewStory(v.Id, v.Username, v.PublishedOn, *t, *m)
	}

	return fs, nil
}

func (u DownloadUserStories) downloadAndStat(url string) (*storage.File, error) {
	f, err := u.client.Download(http.GetHttpRequest(url, nil, nil))

	if err != nil {
		return nil, err
	}

	s, err := f.Stat()

	if err != nil {
		return nil, err
	}

	return &storage.File{
		ModifyDateTime: s.ModTime(),
		IsDir:          s.IsDir(),
		Path:           f.Name(),
	}, err
}
