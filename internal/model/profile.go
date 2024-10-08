package model

type Profile struct {
	Bio     Bio
	Stories []Story
}

func NewProfile(
	bio Bio,
	stories []Story,
) Profile {
	return Profile{
		Bio:     bio,
		Stories: stories,
	}
}
