package archive

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Implements [ArchiveRepository] using [FileSystemStorage] as the data source.
type FileSystemArchiveRepository struct {
	storage           *storage.FileSystemStorage
	client            http.HttpClient
	virtualContentDir string
}

func (r FileSystemArchiveRepository) All() ([]string, error) {
	fs, err := r.storage.Lookup("")

	if err != nil {
		return nil, err
	}

	return model.Map(
		fs,
		func(f storage.File) string {
			return f.Name()
		},
	), nil
}

func (r FileSystemArchiveRepository) Stories(username string) ([]model.Story, error) {
	fs, err := r.storage.Lookup(username)

	if err != nil {
		return nil, err
	}

	return model.MapFilter(
		fs,
		func(f storage.File) (model.Story, error) {
			var s model.Story

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

func (r FileSystemArchiveRepository) Archive(profile model.Profile) (model.Profile, error) {
	tmp := os.TempDir()
	alreadyDownloadedContent, err := r.storage.Lookup(profile.Bio.Username)

	if err != nil {
		alreadyDownloadedContent = []storage.File{}
	}

	logging.LogInfo("downloading stories of %s...", profile.Bio.Username)
	stories, err := r.downloadStories(profile, alreadyDownloadedContent, tmp)
	if err != nil {
		return model.Profile{}, err
	}

	logging.LogInfo("downloading avatar of %s...", profile.Bio.Username)
	avatar, err := r.downloadAvatar(profile, alreadyDownloadedContent, tmp)
	if err != nil {
		return model.Profile{}, err
	}

	return model.NewProfile(
		model.NewBio(
			profile.Bio.Username, profile.Bio.Description, avatar, profile.Bio.IsPrivate),
		stories,
	), nil
}

func (r FileSystemArchiveRepository) downloadStories(profile model.Profile, alreadyDownloadedContent []storage.File, path string) ([]model.Story, error) {
	stories := profile.Stories

	for i, v := range stories {
		// check if content has been downloaded already
		if f := model.Filter(alreadyDownloadedContent, func(f storage.File) bool {
			return strings.Contains(f.Name(), v.Id)
		}); len(f) == 3 {
			dat := model.Find(f, func(f storage.File) bool {
				return strings.HasSuffix(f.Path, ".dat")
			})

			if dat == nil {
				logging.LogError("skipping archived story (%s). missing .dat file", v.Id)
				continue
			}

			var s model.Story

			bs, err := os.ReadFile(dat.Path)

			if err != nil {
				logging.LogError("skipping archived story (%s). failed reading .dat file", v.Id)
				continue
			}

			err = json.Unmarshal(bs, &s)

			if err != nil {
				logging.LogError("skipping archived story (%s). failed parsing .dat file", v.Id)
				continue
			}

			stories[i] = s
			continue
		}

		t, terr := r.downloadAndStat(v.Thumbnail, fmt.Sprintf("%s/%s.jpeg", path, v.Id))

		if terr != nil {
			return nil, terr
		}

		murl := fmt.Sprintf("%s/%s.mp4", path, v.Id)

		if !v.IsVideo {
			murl = fmt.Sprintf("%s/%s-media.jpeg", path, v.Id)
		}

		m, merr := r.downloadAndStat(v.Media, murl)
		if merr != nil {
			return nil, merr
		}

		fs, err := r.storage.Store(profile.Bio.Username, []storage.File{*t, *m})

		if err != nil {
			return nil, terr
		}

		s := model.NewStory(
			v.Id,
			v.Username,
			v.PublishedOn,
			v.IsVideo,
			r.purifyUrl(fs[0].Path),
			r.purifyUrl(fs[1].Path),
		)

		bs, err := json.Marshal(s)

		if err != nil {
			return stories, err
		}

		_, err = r.storage.StoreRaw(fmt.Sprintf("%s/%s.dat", v.Username, v.Id), bs)

		if err != nil {
			return stories, err
		}

		stories[i] = s
	}

	return stories, nil
}

func (r FileSystemArchiveRepository) downloadAvatar(profile model.Profile, alreadyDownloadedContent []storage.File, path string) (string, error) {
	if f := model.Find(alreadyDownloadedContent, func(f storage.File) bool {
		return f.Name() == "avatar.jpeg"
	}); f != nil {
		return r.purifyUrl(f.Path), nil
	}

	af, err := r.downloadAndStat(profile.Bio.Avatar, fmt.Sprintf("%s/avatar.jpeg", path))

	if err != nil {
		return "", err
	}

	fs, err := r.storage.Store(profile.Bio.Username, []storage.File{*af})

	if err != nil {
		return "", err
	}

	return r.purifyUrl(fs[0].Path), nil
}

func (r FileSystemArchiveRepository) purifyUrl(url string) string {
	if strings.HasPrefix(url, r.storage.Root) {
		return strings.Replace(url, r.storage.Root, r.virtualContentDir, 1)
	}

	return url
}

func (r FileSystemArchiveRepository) downloadAndStat(url string, destPath string) (*storage.File, error) {
	logging.LogInfo("downloading %s...", url)
	f, err := r.client.Download(http.GetHttpRequest(url, nil, nil), destPath)

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
