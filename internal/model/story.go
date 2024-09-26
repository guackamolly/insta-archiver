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

// A story model that can be used for client responses.
type CloudStory = Story[string]

// A story model that lives within the host machine.
type FileStory = Story[storage.File]
