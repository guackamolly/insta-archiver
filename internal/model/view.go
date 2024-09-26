package model

import (
	"time"
)

type ArchivedUserView struct {
	Username        string
	Description     string
	LastStories     []CloudStory
	ArchivedStories map[string][]CloudStory
}

func NewArchivedUserView(
	username,
	description string,
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
		Username:        username,
		Description:     description,
		LastStories:     ts,
		ArchivedStories: as,
	}
}
