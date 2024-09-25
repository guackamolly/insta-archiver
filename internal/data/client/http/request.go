package http

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
)

type Headers map[string]string
type QueryParameters map[string]string

// Abstracts the definition of an HTTP request.
type HttpRequest struct {
	URL       string
	Verb      string
	MediaType string
	Cache     string
	Body      []byte
	Headers   Headers
	Query     QueryParameters
}

// Creates an [HttpRequest] suitable for GET requests.
func GetHttpRequest(
	url string,
	headers *Headers,
	queryParams *QueryParameters,
) HttpRequest {
	hds := headers
	qps := queryParams

	if hds == nil {
		hds = &Headers{}
	}

	if qps == nil {
		qps = &QueryParameters{}
	}

	return HttpRequest{
		Verb:    "GET",
		URL:     url,
		Headers: *hds,
		Query:   *qps,
	}
}

// Creates an [HttpRequest] suitable for POST requests.
func PostHttpRequest(
	url string,
	mediaType *string,
	body *[]byte,
	headers *Headers,
	queryParams *QueryParameters,
) HttpRequest {
	var mt string
	b := body
	hds := headers
	qps := queryParams

	if mediaType == nil {
		mt = "application/json"
	} else {
		mt = *mediaType
	}

	if b == nil {
		b = &[]byte{}
	}

	if hds == nil {
		hds = &Headers{}
	}

	if qps == nil {
		qps = &QueryParameters{}
	}

	return HttpRequest{
		Verb:      "POST",
		URL:       url,
		MediaType: mt,
		Body:      *b,
		Headers:   *hds,
		Query:     *qps,
	}
}

// Converts the structure into a native [http.Request] structure.
func (r HttpRequest) Native() (*http.Request, error) {
	req, err := http.NewRequest(r.Verb, r.QueryURL(), bytes.NewBuffer(r.Body))

	if err != nil {
		return req, err
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	if len(r.Cache) > 0 {
		req.Header.Set("cache-control", string(r.Cache))
	}

	return req, err
}

// Combines URL and QueryParameters in a single string (e.g., https://www.example.com?lang=en&page=2)
func (r HttpRequest) QueryURL() string {
	u := r.URL

	if len(r.Query) > 0 {
		var buf strings.Builder

		buf.WriteString(u)
		buf.WriteByte('?')

		for k, v := range r.Query {
			buf.WriteByte('&')
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}

		u = buf.String()
	}

	return u
}
