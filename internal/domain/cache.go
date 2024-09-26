package domain

import (
	"fmt"

	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type LoadCacheArchivedUserView struct {
	repository cache.CacheRepository[string, model.ArchivedUserView]
}

type CacheArchivedUserView struct {
	repository cache.CacheRepository[string, model.ArchivedUserView]
}

type GetCachedArchivedUserView struct {
	repository cache.CacheRepository[string, model.ArchivedUserView]
}

func (u LoadCacheArchivedUserView) Invoke() error {
	c, err := u.repository.Load()

	for id := range c {
		fmt.Printf("loaded cache for user %s\n", id)
	}

	return err
}

func (u CacheArchivedUserView) Invoke(view model.ArchivedUserView) (model.ArchivedUserView, error) {
	var v model.ArchivedUserView
	ce, err := u.repository.Update(view.Username, view)

	if err == nil {
		v = ce.Value
	}

	return v, err
}

func (u GetCachedArchivedUserView) Invoke(username string) (*model.ArchivedUserView, error) {
	v, err := u.repository.Lookup(username)

	if err == nil {
		return &v.Value, nil
	}

	return nil, nil
}
