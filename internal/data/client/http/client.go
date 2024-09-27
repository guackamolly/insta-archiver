package http

import "os"

// Abstracts the behaviour and usage of the application HTTP Client.
type HttpClient interface {
	Do(req HttpRequest) (HttpResponse, error)
	Download(req HttpRequest, destPath string) (*os.File, error)
}
