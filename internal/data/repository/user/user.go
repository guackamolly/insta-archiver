package user

import (
	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Operations the application is able to do in regards to an Instagram user.
type UserRepository interface {
	Stories(username string) ([]model.CloudStory, error)
}

func NewViewIGStoryUserRepository(
	client http.HttpClient,
) UserRepository {
	return ViewIGStoryUserRepository{
		client: client,
	}
}

func NewAnonyIGStoryUserRepository(
	client http.HttpClient,
) UserRepository {
	return AnonyIGStoryUserRepository{
		client: client,
	}
}
