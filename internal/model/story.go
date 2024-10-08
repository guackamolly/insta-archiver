package model

import (
	"time"
)

type Story struct {
	Id          string
	Username    string
	PublishedOn time.Time
	Thumbnail   string
	Media       string
	IsVideo     bool
}

func NewStory(
	id,
	username string,
	publishedOn time.Time,
	isVideo bool,
	thumbnail string,
	media string,
) Story {
	return Story{
		Id:          id,
		Username:    username,
		PublishedOn: publishedOn,
		Thumbnail:   thumbnail,
		Media:       media,
		IsVideo:     isVideo,
	}
}
