package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Implements [UserRepository] using [ViewIGStory](https://viewigstory.com) API.
type ViewIGStoryUserRepository struct {
	client http.HttpClient
}

func (r ViewIGStoryUserRepository) Stories(username string) ([]model.Story, error) {
	var res []model.Story

	u := fmt.Sprintf("https://viewigstory.com/api/stories/%s", username)
	resp, err := r.client.Do(
		http.PostHttpRequest(
			u,
			nil,
			nil,
			headers(u),
			nil,
		),
	)

	if err != nil {
		return res, err
	}

	stories, err := http.Typed[getUserStoriesResponse](&resp.Body)

	if err != nil {
		return res, err
	}

	res = make([]model.Story, len(stories.LastStories))
	for i, v := range stories.LastStories {
		pdt, err := time.Parse(time.DateOnly, v.CreatedTime)

		if err != nil {
			pdt = time.Now()
		}

		res[i] = model.NewStory(
			v.CreatedTime,
			username,
			pdt,
			v.Type == "video",
			fmt.Sprintf("https://viewigstory.com/proxy/%s", v.ThumbnailURL),
			fmt.Sprintf("https://viewigstory.com/proxy/%s", v.VideoURL),
		)
	}

	return res, err
}

func (r ViewIGStoryUserRepository) Bio(username string) (model.Bio, error) {
	return model.Bio{}, errors.New("not implemented yet")
}

// Structure of a successful JSON response from ViewIGStory API.
type getUserStoriesResponse struct {
	LastStories []struct {
		CreatedTime  string `json:"createdTime"`
		Type         string `json:"type"`
		ThumbnailURL string `json:"thumbnailUrl"`
		VideoURL     string `json:"videoUrl"`
	} `json:"lastStories"`
}

// Custom headers extracted from a browser web request. Required to not get flagged in ViewIGStory API.
func headers(url string) *http.Headers {
	return &http.Headers{
		"accept":             "*/*",
		"accept-language":    "en-US,en;q=0.9",
		"priority":           "u=1, i",
		"sec-ch-ua":          `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Linux"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"Referrer-Policy":    "strict-origin-when-cross-origin",
		"User-Agent":         "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Referer":            url,
	}
}
