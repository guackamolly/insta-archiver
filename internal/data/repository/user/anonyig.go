package user

import (
	"fmt"
	"time"

	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Implements [UserRepository] using [Anonyig](https://anonyig.com/) API.
type AnonyIGStoryUserRepository struct {
	client http.HttpClient
}

func (r AnonyIGStoryUserRepository) Stories(username string) ([]model.CloudStory, error) {
	var res []model.CloudStory

	u := "https://anonyig.com/api/ig/story"
	resp, err := r.client.Do(
		http.GetHttpRequest(
			u,
			anonyigHeaders(u),
			&http.QueryParameters{"url": fmt.Sprintf("https://www.instagram.com/stories/%s/", username)},
		),
	)

	if err != nil {
		return res, err
	}

	body, err := http.Typed[AnonyIGGetUserStoriesResponse](&resp.Body)

	if err != nil {
		return res, err
	}

	stories := body.Result

	res = make([]model.CloudStory, len(stories))
	for i, v := range stories {
		if len(v.ImageVersions2.Candidates) == 0 {
			continue
		}

		thumb := model.Find(v.ImageVersions2.Candidates, func(c anonyIGGetUserStoriesResponseThumbnailCandidate) bool {
			// 1136 height is the best candidate
			return c.Height == 1136
		})

		video := model.Find(v.VideoVersions, func(c anonyIGGetUserStoriesResponseVideoCandidate) bool {
			// 102 is low res one
			return c.Type == 102
		})

		if thumb == nil {
			thumb = &v.ImageVersions2.Candidates[0]
		}

		thumbnail := thumb.URL
		media := thumbnail

		if video != nil {
			media = video.URL
		}

		pdt, err := time.Parse(time.StampMilli, fmt.Sprintf("%d000", v.TakenAt))

		if err != nil {
			pdt = time.Now()
		} else {
			pdt = pdt.Add(time.Hour * 24)
		}

		res[i] = model.NewStory(
			v.Pk,
			username,
			pdt,
			thumbnail,
			media,
		)
	}

	return res, err
}

type anonyIGGetUserStoriesResponseThumbnailCandidate struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	URL          string `json:"url"`
	URLSignature struct {
		Expires   string `json:"expires"`
		Signature string `json:"signature"`
	} `json:"url_signature"`
}

type anonyIGGetUserStoriesResponseVideoCandidate struct {
	Height       int    `json:"height"`
	Type         int    `json:"type"`
	URL          string `json:"url"`
	Width        int    `json:"width"`
	URLSignature struct {
		Expires   string `json:"expires"`
		Signature string `json:"signature"`
	} `json:"url_signature"`
}

// Structure of a successful JSON response from ViewIGStory API.
type AnonyIGGetUserStoriesResponse struct {
	Result []struct {
		ImageVersions2 struct {
			Candidates []anonyIGGetUserStoriesResponseThumbnailCandidate `json:"candidates"`
		} `json:"image_versions2"`
		Pk            string                                        `json:"pk"`
		TakenAt       int                                           `json:"taken_at"`
		VideoVersions []anonyIGGetUserStoriesResponseVideoCandidate `json:"video_versions"`
	} `json:"result"`
}

// Custom headers extracted from a browser web request. Required to not get flagged in ViewIGStory API.
func anonyigHeaders(url string) *http.Headers {
	return &http.Headers{
		"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.5",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
		"Referer":         url,
	}
}