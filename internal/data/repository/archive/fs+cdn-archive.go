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

// Implements [ArchiveRepository] using both [FileSystemStorage] and a CDN as data sources.
// The CDN will store stories, while [FileSystemStorage] stores .dat files.
type FileSystemCDNArchiveRepository struct {
	storage *storage.FileSystemStorage
	client  http.HttpClient
	cdnUrl  string
}

func (r FileSystemCDNArchiveRepository) All() ([]string, error) {
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

func (r FileSystemCDNArchiveRepository) Stories(username string) ([]model.Story, error) {
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

func (r FileSystemCDNArchiveRepository) Archive(profile model.Profile) (model.Profile, error) {
	alreadyDownloadedContent, err := r.storage.Lookup(profile.Bio.Username)

	if err != nil {
		alreadyDownloadedContent = []storage.File{}
	}

	logging.LogInfo("uploading stories of %s...", profile.Bio.Username)
	stories, err := r.uploadStories(profile, alreadyDownloadedContent)
	if err != nil {
		return model.Profile{}, err
	}

	logging.LogInfo("uploading avatar of %s...", profile.Bio.Username)
	avatar, err := r.uploadAvatar(profile)
	if err != nil {
		return model.Profile{}, err
	}

	return model.NewProfile(
		model.NewBio(
			profile.Bio.Username, profile.Bio.Description, avatar, profile.Bio.IsPrivate),
		stories,
	), nil
}

func (r FileSystemCDNArchiveRepository) uploadStories(profile model.Profile, alreadyDownloadedContent []storage.File) ([]model.Story, error) {
	stories := profile.Stories

	for i, v := range stories {
		// check if content has been downloaded already
		if f := model.Filter(alreadyDownloadedContent, func(f storage.File) bool {
			return strings.Contains(f.Name(), v.Id)
		}); len(f) > 0 {
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

		t, err := r.uploadCDN(profile.Bio.Username, v.Thumbnail, false)
		if err != nil {
			return nil, err
		}

		m, err := r.uploadCDN(profile.Bio.Username, v.Media, false)
		if err != nil {
			return nil, err
		}

		s := model.NewStory(
			v.Id,
			v.Username,
			v.PublishedOn,
			v.IsVideo,
			t,
			m,
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

func (r FileSystemCDNArchiveRepository) uploadAvatar(profile model.Profile) (string, error) {
	return r.uploadCDN(profile.Bio.Username, profile.Bio.Avatar, true)
}

func (r FileSystemCDNArchiveRepository) uploadCDN(
	username string,
	url string,
	override bool,
) (string, error) {
	he := http.Headers{
		"x-download-url": url,
		"x-dir":          username,
		"x-override":     fmt.Sprintf("%v", override),
	}
	req := http.PostHttpRequest(r.cdnUrl, nil, nil, &he, nil)

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.Nok() {
		return "", fmt.Errorf("failed to upload file (%s) to CDN, %d", url, resp.StatusCode)
	}

	return resp.RequestURL, nil
}
