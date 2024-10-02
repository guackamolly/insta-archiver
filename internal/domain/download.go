package domain

import (
	"fmt"
	"os"
	"strings"

	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type DownloadUserProfile struct {
	client         http.HttpClient
	contentStorage *storage.FileSystemStorage

	physicalContentDir string
	virtualContentDir  string
}

func (u DownloadUserProfile) Invoke(profile model.Profile) (model.Profile, error) {
	tmp := os.TempDir()
	alreadyDownloadedContent, err := u.contentStorage.Lookup(profile.Bio.Username)

	if err != nil {
		alreadyDownloadedContent = []storage.File{}
	}

	fmt.Printf("downloading stories of %s...\n", profile.Bio.Username)
	stories, err := u.downloadStories(profile, alreadyDownloadedContent, tmp)
	if err != nil {
		return model.Profile{}, err
	}

	fmt.Printf("downloading avatar of %s...\n", profile.Bio.Username)
	avatar, err := u.downloadAvatar(profile, alreadyDownloadedContent, tmp)
	if err != nil {
		return model.Profile{}, err
	}

	return model.NewProfile(
		model.NewBio(
			profile.Bio.Username,
			profile.Bio.Description,
			avatar,
		),
		stories,
	), nil
}

func (u DownloadUserProfile) downloadStories(profile model.Profile, alreadyDownloadedContent []storage.File, path string) ([]model.Story[string], error) {
	stories := profile.Stories

	for i, v := range stories {
		// check if content has been downloaded already
		if f := model.Filter(alreadyDownloadedContent, func(f storage.File) bool {
			return strings.Contains(f.Name(), v.Id)
		}); len(f) == 2 {
			tf := f[0]
			mf := f[1]

			if n := f[0].Name(); strings.HasSuffix(n, ".mp4") || strings.Contains(n, "media.jpeg") {
				tf = f[1]
				mf = f[0]
			}

			stories[i] = model.NewStory(
				v.Id,
				v.Username,
				v.PublishedOn,
				v.IsVideo,
				u.purifyUrl(tf.Path),
				u.purifyUrl(mf.Path),
			)

			continue
		}

		t, terr := u.downloadAndStat(v.Thumbnail, fmt.Sprintf("%s/%s.jpeg", path, v.Id))

		if terr != nil {
			return nil, model.Wrap(terr, DownloadThumbnailFailed)
		}

		murl := fmt.Sprintf("%s/%s.mp4", path, v.Id)

		if !v.IsVideo {
			murl = fmt.Sprintf("%s/%s-media.jpeg", path, v.Id)
		}

		m, merr := u.downloadAndStat(v.Media, murl)
		if merr != nil {
			return nil, model.Wrap(merr, DownloadMediaFailed)
		}

		fs, err := u.contentStorage.Store(profile.Bio.Username, []storage.File{*t, *m})

		if err != nil {
			return nil, model.Wrap(terr, StoreStoriesFailed)
		}

		stories[i] = model.NewStory(
			v.Id,
			v.Username,
			v.PublishedOn,
			v.IsVideo,
			u.purifyUrl(fs[0].Path),
			u.purifyUrl(fs[1].Path),
		)
	}

	return stories, nil
}

func (u DownloadUserProfile) downloadAvatar(profile model.Profile, alreadyDownloadedContent []storage.File, path string) (string, error) {
	if f := model.Find(alreadyDownloadedContent, func(f storage.File) bool {
		return f.Name() == "avatar.jpeg"
	}); f != nil {

		return u.purifyUrl(f.Path), nil
	}

	af, err := u.downloadAndStat(profile.Bio.Avatar, fmt.Sprintf("%s/avatar.jpeg", path))

	if err != nil {
		return "", model.Wrap(err, DownloadAvatarFailed)
	}

	fs, err := u.contentStorage.Store(profile.Bio.Username, []storage.File{*af})

	if err != nil {
		return "", model.Wrap(err, StoreAvatarFailed)
	}

	return u.purifyUrl(fs[0].Path), nil
}

func (u DownloadUserProfile) purifyUrl(url string) string {
	if strings.HasPrefix(url, u.physicalContentDir) {
		return strings.Replace(url, u.physicalContentDir, u.virtualContentDir, 1)
	}

	return url
}

func (u DownloadUserProfile) downloadAndStat(url string, destPath string) (*storage.File, error) {
	fmt.Printf("downloading %s...\n", url)
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
