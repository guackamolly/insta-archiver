package domain

import (
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/model"
)

type CacheArchivedUserView struct {
	storage *storage.MemoryStorage[string, model.ArchivedUserView]
}

type GetCachedArchivedUserView struct {
	storage *storage.MemoryStorage[string, model.ArchivedUserView]
}

func (u CacheArchivedUserView) Invoke(view model.ArchivedUserView) (model.ArchivedUserView, error) {
	return u.storage.Store(view.Username, view)
}

func (u GetCachedArchivedUserView) Invoke(username string) (*model.ArchivedUserView, error) {
	v, err := u.storage.Lookup(username)

	if err == nil {
		return &v, nil
	}

	return nil, nil
}
