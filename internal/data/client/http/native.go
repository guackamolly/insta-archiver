package http

import (
	"net/http"
)

type nativeHttpClient struct {
	client *http.Client
}

// Creates a [HttpClient] that uses native [http] as the http client.
func Native() HttpClient {
	return nativeHttpClient{
		client: http.DefaultClient,
	}
}

func (c nativeHttpClient) Do(req HttpRequest) (HttpResponse, error) {
	var resp HttpResponse

	nreq, err := req.Native()

	if err != nil {
		return resp, err
	}

	nresp, err := c.client.Do(nreq)

	if err != nil {
		return resp, err
	}

	resp = HttpResponse{
		MediaType:  nresp.Header.Get("content-type"),
		StatusCode: nresp.StatusCode,
		Headers:    nresp.Header,
		Body:       HttpBody{read: nresp.Body},
	}

	return resp, err
}
