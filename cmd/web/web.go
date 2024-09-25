package main

import (
	"github.com/guackamolly/insta-archiver/internal/core"
	_http "github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/user"
	"github.com/guackamolly/insta-archiver/internal/http"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize di container
	vault := core.Vault{
		UserRepository: user.ViewIGStoryUserRepository{
			Client: _http.Native(),
		},
	}

	// Initialize echo framework
	e := echo.New()
	defer e.Close()

	// Register server dependencies
	http.RegisterHandlers(e)
	http.RegisterMiddlewares(e, vault)
	http.RegisterStaticFiles(e)
	http.RegisterTemplates(e)

	// Start!
	e.Logger.Fatal(http.Start(e))
}
