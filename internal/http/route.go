package http

var (
	rootRoute    = WithVirtualHost("/")
	archiveRoute = WithVirtualHost("/archive/:id")
)
