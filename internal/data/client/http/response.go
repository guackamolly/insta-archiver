package http

import "net/http"

// Abstracts the definition of an HTTP response.
type HttpResponse struct {
	MediaType  string
	StatusCode int
	Body       HttpBody
	Headers    http.Header
	RequestURL string
}

func (r HttpResponse) Ok() bool {
	return r.StatusCode < 399
}

func (r HttpResponse) Nok() bool {
	return r.StatusCode > 399
}

func (r HttpResponse) Redirection() bool {
	return r.StatusCode > 299 && r.StatusCode < 400
}
