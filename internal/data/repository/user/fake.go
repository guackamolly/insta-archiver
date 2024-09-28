package user

import (
	"time"

	"github.com/guackamolly/insta-archiver/internal/model"
)

// Dumb [UserRepository] that uses fake data for testing purposes.
type FakeUserRepository struct{}

func (r FakeUserRepository) Bio(username string) (model.Bio, error) {
	return model.NewBio(
		username,
		"I have a really cool bio",
		"https://fakeimg.pl/128x128/282828/eae0d0/?retina=1&text=%20%F0%9F%98%8B",
	), nil
}

func (r FakeUserRepository) Stories(username string) ([]model.CloudStory, error) {
	return []model.CloudStory{
		model.NewStory(
			"00001",
			username,
			time.Now(),
			false,
			"https://fakeimg.pl/640x1124/ffffff?text=A",
			"https://fakeimg.pl/640x1124/ffffff?text=A",
		),
		model.NewStory(
			"00002",
			username,
			time.Now(),
			false,
			"https://fakeimg.pl/640x1124/ffffff?text=B",
			"https://fakeimg.pl/640x1124/ffffff?text=B",
		),
		model.NewStory(
			"00003",
			username,
			time.Now(),
			false,
			"https://fakeimg.pl/640x1124/ffffff?text=C",
			"https://fakeimg.pl/640x1124/ffffff?text=C",
		),
		model.NewStory(
			"00004",
			username,
			time.Now(),
			false,
			"https://fakeimg.pl/640x1124/ffffff?text=D",
			"https://fakeimg.pl/640x1124/ffffff?text=D",
		),
		model.NewStory(
			"00005",
			username,
			time.Now(),
			false,
			"https://fakeimg.pl/640x1124/ffffff?text=E",
			"https://fakeimg.pl/640x1124/ffffff?text=E",
		),
	}, nil
}
