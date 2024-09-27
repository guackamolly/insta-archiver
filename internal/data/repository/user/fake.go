package user

import (
	"time"

	"github.com/guackamolly/insta-archiver/internal/model"
)

// Dumb [UserRepository] that uses fake data for testing purposes.
type FakeUserRepository struct{}

func (r FakeUserRepository) Stories(username string) ([]model.CloudStory, error) {
	return []model.CloudStory{
		model.NewStory(
			"00001",
			username,
			time.Now(),
			"https://fakeimg.pl/640x1124/ffffff?text=A",
			"https://fakeimg.pl/640x1124/ffffff?text=A",
		),
		model.NewStory(
			"00002",
			username,
			time.Now(),
			"https://fakeimg.pl/640x1124/ffffff?text=B",
			"https://fakeimg.pl/640x1124/ffffff?text=B",
		),
		model.NewStory(
			"00003",
			username,
			time.Now(),
			"https://fakeimg.pl/640x1124/ffffff?text=C",
			"https://fakeimg.pl/640x1124/ffffff?text=C",
		),
		model.NewStory(
			"00004",
			username,
			time.Now(),
			"https://fakeimg.pl/640x1124/ffffff?text=D",
			"https://fakeimg.pl/640x1124/ffffff?text=D",
		),
		model.NewStory(
			"00005",
			username,
			time.Now(),
			"https://fakeimg.pl/640x1124/ffffff?text=E",
			"https://fakeimg.pl/640x1124/ffffff?text=E",
		),
	}, nil
}
