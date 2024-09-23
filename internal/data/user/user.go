package user

import (
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Operations the application is able to do in regards to an Instagram user.
type UserRepository interface {
	Stories(uid string) ([]model.Story, error)
}
