package model

type Profile struct {
	Bio     Bio
	Stories []Story[string]
}

func NewProfile(
	bio Bio,
	stories []Story[string],
) Profile {
	return Profile{
		Bio:     bio,
		Stories: stories,
	}
}
