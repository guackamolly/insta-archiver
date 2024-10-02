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
	c, err := WrapResult0(u.repository.Load, LoadCacheFailed)

	for id := range c {
		fmt.Printf("loaded cache for user %s\n", id)
	}

	return err
}

func (u CacheArchivedUserView) Invoke(view model.ArchivedUserView) (model.ArchivedUserView, error) {
	ce, err := u.repository.Update(view.Username, view)

	if err == nil {
		return ce.Value, nil
	}

	return view, model.Wrap(err, UpdateCacheFailed)
}

func (u GetCachedArchivedUserView) Invoke(username string) (model.ArchivedUserView, error) {
	v, err := u.repository.Lookup(username)

	if err == nil && !v.IsOutdated() {
		return v.Value, nil
	}

	return model.ArchivedUserView{}, model.Wrap(err, LookupCacheFailed)
}
