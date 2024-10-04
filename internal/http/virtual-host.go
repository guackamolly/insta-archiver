package http

import (
	"fmt"
	"os"
	"strings"
)

const (
	serverVirtualHostEnvKey = "server_virtual_host"
)

var (
	serverVirtualHost = os.Getenv(serverVirtualHostEnvKey)
)

func init() {
	if !strings.HasPrefix(serverVirtualHost, "/") {
		serverVirtualHost = "/" + serverVirtualHost
	}

	if !strings.HasSuffix(serverVirtualHost, "/") {
		serverVirtualHost = serverVirtualHost + "/"
	}
}

func WithVirtualHost(path string) string {
	if path == "" || path == "/" {
		return serverVirtualHost
	}

	if serverVirtualHost == "/" {
		return path
	}

	return fmt.Sprintf("%s%s", serverVirtualHost, path)
}
