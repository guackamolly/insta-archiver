package domain

import (
	"time"

	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
)

var cacheRefreshTimers = map[string]<-chan time.Time{}

type LoadCacheArchivedUserView struct {
	repository cache.CacheRepository[string, model.ArchivedUserView]
}

type CacheArchivedUserView struct {
	cacheRepo   cache.CacheRepository[string, model.ArchivedUserView]
	archiveRepo archive.ArchiveRepository
}

type GetCachedArchivedUserView struct {
	repository cache.CacheRepository[string, model.ArchivedUserView]
}

func (u LoadCacheArchivedUserView) Invoke() error {
	c, err := WrapResult0(u.repository.Load, LoadCacheFailed)

	for id := range c {
		logging.LogInfo("loaded cache for user %s", id)

	}

	return err
}

func (u CacheArchivedUserView) Invoke(view model.ArchivedUserView) (model.ArchivedUserView, error) {
	ce, err := u.cacheRepo.Update(view.Username, view)

	if err == nil {
		return ce.Value, nil
	}

	return view, model.Wrap(err, UpdateCacheFailed)
}

func (u CacheArchivedUserView) Schedule(username string, refresh func() (model.ArchivedUserView, error)) error {
	if _, ok := cacheRefreshTimers[username]; ok {
		logging.LogWarning("%s is already scheduled for refresh. skipping...", username)
		return nil
	}

	duration := u.cacheRepo.Policy()

	t := time.Tick(duration)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				delete(cacheRefreshTimers, username)
				go func() {
					u.Schedule(username, refresh)
				}()
			}
		}()
		for range t {
			view, err := refresh()

			if err != nil {
				logging.LogError("failed to obtain view to refresh...%v", err)
				continue
			}

			view, err = u.Invoke(view)

			if err != nil {
				logging.LogError("failed to refresh cache...%v", err)
			} else {
				logging.LogInfo("cache of %s view has been refreshed", view.Username)
			}
		}
	}()

	cacheRefreshTimers[username] = t
	logging.LogInfo("scheduled cache refresh for user %s with duration %s", username, duration)

	return nil
}

func (u CacheArchivedUserView) ScheduleAll(
	refresh func(username string) (model.ArchivedUserView, error),
) error {
	usernames, err := u.archiveRepo.All()
	if err != nil {
		return model.Wrap(err, ScheduleArchiveFailed)
	}

	for _, un := range usernames {
		err = u.Schedule(un, func() (model.ArchivedUserView, error) {
			return refresh(un)
		})

		if err != nil {
			logging.LogError("failed to schedule archive for user %s. skipping until next hit", un)
		}
	}

	return nil
}

func (u GetCachedArchivedUserView) Invoke(username string) (model.ArchivedUserView, error) {
	v, err := u.repository.Lookup(username)

	if err == nil && !v.IsOutdated() {
		return v.Value, nil
	}

	return model.ArchivedUserView{}, model.Wrap(err, LookupCacheFailed)
}
