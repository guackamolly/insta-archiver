package model

import (
	"time"

	"github.com/guackamolly/insta-archiver/internal/data/storage"
)

// Central model of the application. Defines the required data for an Instagram Story.
// Typed so it's possible to differentiate a story accessible by clients (string) or accessible in the host machine (os.File)
type Story[T any] struct {
	Id          string
	Username    string
	PublishedOn time.Time
	Thumbnail   T
	Media       T
}

func NewStory[T any](
	id,
	username string,
	publishedOn time.Time,
	thumbnail T,
	media T,
) Story[T] {
	return Story[T]{
		Id:          id,
		Username:    username,
		PublishedOn: publishedOn,
		Thumbnail:   thumbnail,
		Media:       media,
	}
}

func (s Story[T]) Equal(o Story[T]) bool {
	return s.Username == o.Username && (s.Id == o.Id || s.PublishedOn == o.PublishedOn)
}

// A story model that can be used for client responses.
type CloudStory = Story[string]

// A story model that lives within the host machine.
type FileStory = Story[storage.File]
