package user

import (
	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Operations the application is able to do in regards to an Instagram user.
type UserRepository interface {
	Bio(username string) (model.Bio, error)
	Stories(username string) ([]model.Story, error)
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

func NewFakeUserRepository() FakeUserRepository {
	return FakeUserRepository{}
}
