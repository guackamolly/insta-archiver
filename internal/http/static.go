package http

import (
	"os"

	"github.com/labstack/echo/v4"
)

const (
	serverPublicRootEnvKey = "server_public_root"
)

var (
	serverPublicRoot = os.Getenv(serverPublicRootEnvKey)

	files = map[string]string{
		"/index.html":    serverPublicRoot + "index.html",
		"/index.css":     serverPublicRoot + "index.css",
		"/manifest.json": serverPublicRoot + "manifest.json",
	}

	dirs = map[string]string{
		"/static": serverPublicRoot + "static/",
	}

	templates = []string{
		serverPublicRoot + "archive/index.html",
	}

	errors = map[int]string{
		404: serverPublicRoot + "404/index.html",
	}

	root = files["/index.html"]
)

func RegisterStaticFiles(e *echo.Echo) error {
	for k, v := range files {
		e.File(k, v)
	}

	for k, v := range dirs {
		e.Static(k, v)
	}

	return nil
}
