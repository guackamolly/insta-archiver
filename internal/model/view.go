package model

import (
	"time"
)

type ArchivedUserView struct {
	Username             string
	Bio                  Bio
	LastStories          []CloudStory
	ArchivedStories      map[string][]CloudStory
	ArchivedStoriesCount int
}

func NewArchivedUserView(
	username string,
	bio Bio,
	stories []CloudStory,
) ArchivedUserView {
	as := GroupBy(stories, func(s CloudStory) string {
		return s.PublishedOn.Format(time.DateOnly)
	})

	tk := time.Now().Format(time.DateOnly)
	ts, ok := as[tk]

	if ok {
		delete(as, tk)
	}

	return ArchivedUserView{
		Username:             username,
		Bio:                  bio,
		LastStories:          ts,
		ArchivedStories:      as,
		ArchivedStoriesCount: len(stories),
	}
}
