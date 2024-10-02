package domain

import (
	"time"

	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
)

var cacheRefreshTimers = map[string]<-chan time.Time{}

type LoadCacheArchivedUserView struct {
	repository cache.CacheRepository[string, model.ArchivedUserView]
}

type CacheArchivedUserView struct {
	repository cache.CacheRepository[string, model.ArchivedUserView]
}

type GetCachedArchivedUserView struct {
	repository cache.CacheRepository[string, model.ArchivedUserView]
}

func (u LoadCacheArchivedUserView) Invoke() ([]cache.CacheEntry[model.ArchivedUserView], error) {
	c, err := WrapResult0(u.repository.Load, LoadCacheFailed)

	vs := make([]cache.CacheEntry[model.ArchivedUserView], len(c))

	i := 0
	for id, ce := range c {
		logging.LogInfo("loaded cache for user %s", id)
		vs[i] = ce

		i++
	}

	return vs, err
}

func (u CacheArchivedUserView) Invoke(view model.ArchivedUserView) (model.ArchivedUserView, error) {
	ce, err := u.repository.Update(view.Username, view)

	if err == nil {
		return ce.Value, nil
	}

	return view, model.Wrap(err, UpdateCacheFailed)
}

func (u CacheArchivedUserView) Schedule(username string, duration time.Duration, refresh func() (model.ArchivedUserView, error)) error {
	if _, ok := cacheRefreshTimers[username]; ok {
		logging.LogWarning("%s is already scheduled for refresh. skipping...", username)
		return nil
	}

	if duration == 0 {
		ce, err := u.repository.Lookup(username)

		if err != nil {
			return err
		}

		duration = ce.RefreshPolicy
	}

	t := time.Tick(duration)
	go func() {
		for range t {
			view, err := refresh()

			if err != nil {
				logging.LogError("failed to obtain view to refresh...\n%v", err)
				continue
			}

			view, err = u.Invoke(view)

			if err != nil {
				logging.LogError("failed to refresh cache...\n%v", err)
			} else {
				logging.LogInfo("cache of %s view has been refreshed", view.Username)
			}
		}
	}()

	cacheRefreshTimers[username] = t
	logging.LogInfo("scheduled cache refresh with duration %s", duration)

	return nil
}

func (u GetCachedArchivedUserView) Invoke(username string) (model.ArchivedUserView, error) {
	v, err := u.repository.Lookup(username)

	if err == nil && !v.IsOutdated() {
		return v.Value, nil
	}

	return model.ArchivedUserView{}, model.Wrap(err, LookupCacheFailed)
}
