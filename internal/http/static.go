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
		"/static":  serverPublicRoot + "static/",
		"/content": serverPublicRoot + "content/",
	}

	templates = []string{
		serverPublicRoot + "archive/index.html",
	}

	errors = map[int]string{
		404: serverPublicRoot + "404/index.html",
	}

	root     = files["/index.html"]
	fallback = root
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

// Returns a tuples that identifies:
// [0] - Content directory that is accessible through the file system (physical)
// [1] - Content directory that is accessible through the network (virtual)
func ContentDir() [2]string {
	return [2]string{
		dirs["/content"],
		"/content/",
	}
}
