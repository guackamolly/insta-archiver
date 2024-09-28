package model

type Bio struct {
	Username    string
	Description string
	Avatar      string
}

func NewBio(
	username string,
	description string,
	avatar string,
) Bio {
	return Bio{
		Username:    username,
		Description: description,
		Avatar:      avatar,
	}
}

func DefaultBio() Bio {
	return NewBio(
		"404",
		"404",
		"",
	)
}
