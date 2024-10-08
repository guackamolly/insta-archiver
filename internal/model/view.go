package model

import (
	"slices"
	"strings"
	"time"
)

type archiveStoriesItem struct {
	Date    string
	Stories []Story
}

type ArchivedUserView struct {
	Username             string
	Description          string
	Avatar               string
	IsPrivate            bool
	LastStories          []Story
	ArchivedStories      []archiveStoriesItem
	ArchivedStoriesCount int
}

func NewArchivedUserView(
	profile Profile,
) ArchivedUserView {
	as := GroupBy(profile.Stories, func(s Story) string {
		return s.PublishedOn.Format(time.DateOnly)
	})

	tk := time.Now().Format(time.DateOnly)
	ts, ok := as[tk]

	if ok {
		delete(as, tk)
	}

	sas := make([]archiveStoriesItem, len(as))
	i := 0
	for k, v := range as {
		sas[i] = archiveStoriesItem{
			Date:    k,
			Stories: v,
		}
		i++
	}
	slices.SortFunc(sas, func(x archiveStoriesItem, y archiveStoriesItem) int {
		return -strings.Compare(x.Date, y.Date)
	})

	return ArchivedUserView{
		Username:             profile.Bio.Username,
		Description:          profile.Bio.Description,
		Avatar:               profile.Bio.Avatar,
		IsPrivate:            profile.Bio.IsPrivate,
		LastStories:          ts,
		ArchivedStories:      sas,
		ArchivedStoriesCount: len(profile.Stories),
	}
}
