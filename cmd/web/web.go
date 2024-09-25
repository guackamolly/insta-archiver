package main

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	client "github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/repository/archive"
	"github.com/guackamolly/insta-archiver/internal/data/repository/user"
	"github.com/guackamolly/insta-archiver/internal/data/storage"
	"github.com/guackamolly/insta-archiver/internal/domain"
	"github.com/guackamolly/insta-archiver/internal/http"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize echo framework
	e := echo.New()
	defer e.Close()

	// Initialize di container
	vault, err := createVault()

	if err != nil {
		e.Logger.Fatal(err)
	}

	// Register server dependencies
	http.RegisterHandlers(e)
	http.RegisterMiddlewares(e, vault)
	http.RegisterStaticFiles(e)
	http.RegisterTemplates(e)

	// Start!
	e.Logger.Fatal(http.Start(e))
}

func createVault() (core.Vault, error) {
	var vault core.Vault

	client := client.Native()
	fstorage, err := storage.NewFileSystemStorage("static/")

	if err != nil {
		return vault, err
	}

	archiveRepo := archive.NewFileSystemArchiveRepository(fstorage)
	userRepo := user.NewViewIGStoryUserRepository(client)

	vault = core.Vault{
		PurifyUsername:      domain.NewPurifyUsername(),
		DownloadUserStories: domain.NewDownloadUserStories(client),
		GetLatestStories:    domain.NewGetLatestStories(userRepo),
		ArchiveUserStories:  domain.NewArchiveUserStories(archiveRepo),
	}

	return vault, nil
}
