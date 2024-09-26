package model

type ArchivedUserView struct {
	Username        string
	Description     string
	LastStories     []CloudStory
	ArchivedStories []CloudStory
}
