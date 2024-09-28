package main

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	client "github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/domain"
	"github.com/guackamolly/insta-archiver/internal/http"
	"github.com/guackamolly/insta-archiver/internal/model"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize echo framework
	e := echo.New()
	defer e.Close()

	// Initialize di container
	contentDir := http.ContentDir()
	vault, err := createVault(contentDir[0], contentDir[1])

	if err != nil {
		e.Logger.Fatal(err)
	}

	// Register server dependencies
	http.RegisterHandlers(e)
	http.RegisterMiddlewares(e, vault)
	http.RegisterStaticFiles(e)
	http.RegisterTemplates(e)

	// Call beforeStart hook
	http.BeforeStart(e, vault)

	// Start!
	e.Logger.Fatal(http.Start(e))
}

func createVault(
	physicalContentDir,
	virtualContentDir string,
) (core.Vault, error) {
	var vault core.Vault

	client := client.Native()
	fStorage, err := storage.NewFileSystemStorage(physicalContentDir)
	auvMStorage := storage.NewMemoryStorage[string, cache.CacheEntry[model.ArchivedUserView]]()
	bioMStorage := storage.NewMemoryStorage[string, cache.CacheEntry[model.Bio]]()

	if err != nil {
		return vault, err
	}

	archiveRepo := archive.NewFileSystemArchiveRepository(fStorage)
	userRepo := user.NewAnonyIGStoryUserRepository(client)
	auvCacheRepo := cache.NewFileSystemMemoryCacheRepository(fStorage, auvMStorage)
	bioCacheRepo := cache.NewMemoryCacheRepository(bioMStorage)

	vault = core.Vault{
		FilterStoriesForDownload:  domain.NewFilterStoriesForDownload(archiveRepo),
		PurifyStaticUrl:           domain.NewPurifyStaticUrl(physicalContentDir, virtualContentDir),
		PurifyUsername:            domain.NewPurifyUsername(),
		DownloadUserAvatar:        domain.NewDownloadUserAvatar(client),
		DownloadUserStories:       domain.NewDownloadUserStories(client),
		GetArchivedStories:        domain.NewGetArchivedStories(archiveRepo),
		GetLatestStories:          domain.NewGetLatestStories(userRepo),
		ArchiveUserStories:        domain.NewArchiveUserStories(archiveRepo),
		LoadCacheArchivedUserView: domain.NewLoadCacheArchivedUserView(auvCacheRepo),
		ArchiveUserAvatar:         domain.NewArchiveUserAvatar(fStorage),
		CacheArchivedUserView:     domain.NewCacheArchivedUserView(auvCacheRepo),
		GetCachedArchivedUserView: domain.NewGetCachedArchivedUserView(auvCacheRepo),
		GetUserBio:                domain.NewGetUserBio(bioCacheRepo, userRepo),
	}

	return vault, nil
}
