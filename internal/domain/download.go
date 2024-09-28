package domain

import (
	"fmt"
	"os"
	"strings"

	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type DownloadUserStories struct {
	client http.HttpClient
}

type DownloadUserAvatar struct {
	client http.HttpClient
}

type PurifyStaticUrl struct {
	physicalContentDir string
	virtualContentDir  string
}

func (u DownloadUserStories) Invoke(stories []model.CloudStory) ([]model.FileStory, error) {
	fs := make([]model.FileStory, len(stories))
	tmp := os.TempDir()

	for i, v := range stories {
		t, terr := downloadAndStat(u.client, v.Thumbnail, fmt.Sprintf("%s/%s.jpeg", tmp, v.Id))

		if terr != nil {
			return nil, model.Wrap(terr, DownloadThumbnailFailed)
		}

		murl := fmt.Sprintf("%s/%s.mp4", tmp, v.Id)

		if !v.IsVideo {
			murl = fmt.Sprintf("%s/%s-media.jpeg", tmp, v.Id)
		}

		m, merr := downloadAndStat(u.client, v.Media, murl)
		if merr != nil {
			return nil, model.Wrap(merr, DownloadMediaFailed)
		}

		fs[i] = model.NewStory(v.Id, v.Username, v.PublishedOn, v.IsVideo, *t, *m)
	}

	return fs, nil
}

func (u DownloadUserAvatar) Invoke(bio model.Bio) (storage.File, error) {
	var p storage.File
	tmp := os.TempDir()

	ap, err := downloadAndStat(u.client, bio.Avatar, fmt.Sprintf("%s/avatar-%s.jpeg", tmp, bio.Username))

	if err != nil {
		return p, model.Wrap(err, DownloadAvatarFailed)
	}

	return *ap, nil
}

func (u PurifyStaticUrl) Invoke(url string) (string, error) {
	return u.purifyUrl(url), nil
}

func (u PurifyStaticUrl) InvokeStories(stories []model.CloudStory) ([]model.CloudStory, error) {
	cs := make([]model.CloudStory, len(stories))

	for i, s := range stories {
		s.Thumbnail = u.purifyUrl(s.Thumbnail)
		s.Media = u.purifyUrl(s.Media)
		cs[i] = s
	}

	return cs, nil
}

func (u PurifyStaticUrl) purifyUrl(url string) string {
	if strings.HasPrefix(url, u.physicalContentDir) {
		return strings.Replace(url, u.physicalContentDir, u.virtualContentDir, 1)
	}

	return url
}

func downloadAndStat(client http.HttpClient, url string, destPath string) (*storage.File, error) {
	fmt.Printf("downloading %s...\n", url)
	f, err := client.Download(http.GetHttpRequest(url, nil, nil), destPath)

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
