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
	IsVideo     bool
}

func NewStory[T any](
	id,
	username string,
	publishedOn time.Time,
	isVideo bool,
	thumbnail T,
	media T,
) Story[T] {
	return Story[T]{
		Id:          id,
		Username:    username,
		PublishedOn: publishedOn,
		Thumbnail:   thumbnail,
		Media:       media,
		IsVideo:     isVideo,
	}
}

func SameStory[T1 any, T2 any](
	o1 Story[T1],
	o2 Story[T2],
) bool {
	return o1.Username == o2.Username && (o1.Id == o2.Id || o1.PublishedOn == o2.PublishedOn)
}

// A story model that can be used for client responses.
type CloudStory = Story[string]

// A story model that lives within the host machine.
type FileStory = Story[storage.File]
