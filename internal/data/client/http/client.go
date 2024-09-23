package http

// Abstracts the behaviour and usage of the application HTTP Client.
type HttpClient interface {
	Do(req HttpRequest) (HttpResponse, error)
}
