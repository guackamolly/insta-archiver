package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type GetUserBio struct {
	userRepository  user.UserRepository
	cacheRepository cache.MemoryCacheRepository[model.Bio]
}

func (u GetUserBio) Invoke(username string) (model.Bio, error) {
	bio, err := u.userRepository.Bio(username)

	if err == nil {
		return bio, nil
	}

	cache, err := u.cacheRepository.Lookup(username)

	if err == nil {
		return cache.Value, nil
	}

	return bio, model.Wrap(err, FetchBioFailed)

}
