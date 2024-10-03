package model

import (
	"time"
)

type ArchivedUserView struct {
	Username             string
	Description          string
	Avatar               string
	IsPrivate            bool
	LastStories          []Story[string]
	ArchivedStories      map[string][]Story[string]
	ArchivedStoriesCount int
}

func NewArchivedUserView(
	profile Profile,
) ArchivedUserView {
	as := GroupBy(profile.Stories, func(s CloudStory) string {
		return s.PublishedOn.Format(time.DateOnly)
	})

	tk := time.Now().Format(time.DateOnly)
	ts, ok := as[tk]

	if ok {
		delete(as, tk)
	}

	return ArchivedUserView{
		Username:             profile.Bio.Username,
		Description:          profile.Bio.Description,
		Avatar:               profile.Bio.Avatar,
		IsPrivate:            profile.Bio.IsPrivate,
		LastStories:          ts,
		ArchivedStories:      as,
		ArchivedStoriesCount: len(profile.Stories),
	}
}
