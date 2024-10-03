package model

type Bio struct {
	Username    string
	Description string
	Avatar      string
	IsPrivate   bool
}

func NewBio(
	username string,
	description string,
	avatar string,
	isPrivate bool,
) Bio {
	return Bio{
		Username:    username,
		Description: description,
		Avatar:      avatar,
		IsPrivate:   isPrivate,
	}
}

func DefaultBio() Bio {
	return NewBio(
		"404",
		"404",
		"",
		false,
	)
}
