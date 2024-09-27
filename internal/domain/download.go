package domain

import (
	"fmt"
	"os"

	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type DownloadUserStories struct {
	client http.HttpClient
}

func (u DownloadUserStories) Invoke(stories []model.CloudStory) ([]model.FileStory, error) {
	fs := make([]model.FileStory, len(stories))
	tmp := os.TempDir()

	for i, v := range stories {
		t, terr := u.downloadAndStat(v.Thumbnail, fmt.Sprintf("%s/%s.jpeg", tmp, v.Id))

		if terr != nil {
			return nil, model.Wrap(terr, DownloadThumbnailFailed)
		}

		murl := fmt.Sprintf("%s/%s.mp4", tmp, v.Id)

		if !v.IsVideo {
			murl = fmt.Sprintf("%s/%s-media.jpeg", tmp, v.Id)
		}

		m, merr := u.downloadAndStat(v.Media, murl)
		if merr != nil {
			return nil, model.Wrap(merr, DownloadMediaFailed)
		}

		fs[i] = model.NewStory(v.Id, v.Username, v.PublishedOn, v.IsVideo, *t, *m)
	}

	return fs, nil
}

func (u DownloadUserStories) downloadAndStat(url string, destPath string) (*storage.File, error) {
	f, err := u.client.Download(http.GetHttpRequest(url, nil, nil), destPath)

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
