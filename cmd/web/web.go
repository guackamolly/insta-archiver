package main

import (
	"os"

	"github.com/guackamolly/insta-archiver/internal/core"
	client "github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/cache"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/domain"
	"github.com/guackamolly/insta-archiver/internal/http"
	"github.com/guackamolly/insta-archiver/internal/logging"
	"github.com/guackamolly/insta-archiver/internal/model"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize echo framework
	e := echo.New()
	defer e.Close()

	// Initialize logging
	setupLogging(e)

	// Initialize di container
	contentDir := http.ContentDir()
	vault := createVault(contentDir[0], contentDir[1], os.Getenv("cdn_url"))

	// Register server dependencies
	http.RegisterHandlers(e)
	http.RegisterMiddlewares(e, vault)
	http.RegisterStaticFiles(e)
	http.RegisterTemplates(e)

	// Call beforeStart hook
	http.BeforeStart(e, vault)

	// Start!
	logging.LogFatal("server exit %v", http.Start(e))
}

func createVault(
	physicalContentDir,
	virtualContentDir string,
	cdnUrl string,
) core.Vault {
	var vault core.Vault

	client := client.Native()
	contentStorage, err := storage.NewFileSystemStorage(physicalContentDir)

	if err != nil {
		logging.LogError("failed initializing file system storage... %v", err)
	}

	var archiveRepo archive.ArchiveRepository
	if cdnUrl != "" {
		archiveRepo = archive.NewFileSystemCDNArchiveRepository(contentStorage, client, cdnUrl)
	} else {
		archiveRepo = archive.NewFileSystemArchiveRepository(contentStorage, client, virtualContentDir)
	}

	userRepo := user.NewAnonyIGStoryUserRepository(client)
	cacheRepo := cache.NewFileSystemMemoryCacheRepository(contentStorage, storage.NewMemoryStorage[string, cache.CacheEntry[model.ArchivedUserView]]())

	vault = core.Vault{
		PurifyUsername:            domain.NewPurifyUsername(),
		LoadCacheArchivedUserView: domain.NewLoadCacheArchivedUserView(cacheRepo),
		CacheArchivedUserView:     domain.NewCacheArchivedUserView(archiveRepo, cacheRepo),
		GetCachedArchivedUserView: domain.NewGetCachedArchivedUserView(cacheRepo),
		GetUserProfile:            domain.NewGetUserProfile(archiveRepo, userRepo),
		DownloadUserProfile:       domain.NewDownloadUserProfile(archiveRepo),
	}

	return vault
}

func setupLogging(e *echo.Echo) {
	logging.AddLogger(logging.NewConsoleLogger())
	logging.AddLogger(logging.NewEchoLogger(e.Logger))
}
