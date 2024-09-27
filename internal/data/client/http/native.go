package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

func (c nativeHttpClient) Download(req HttpRequest, destPath string) (*os.File, error) {
	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.Nok() {
		return nil, fmt.Errorf("download failed:\nurl: %s\nstatus: %d", req.QueryURL(), resp.StatusCode)
	}

	f, err := c.destOrTempFile(destPath)

	if err == nil {
		defer resp.Body.read.Close()
		_, err = io.Copy(f, resp.Body.read)
	}

	return f, err
}

func (c nativeHttpClient) destOrTempFile(destPath string) (*os.File, error) {
	if destPath == "" {
		return os.CreateTemp("", "*")
	}

	s, err := os.Stat(destPath)

	if os.IsNotExist(err) {
		return os.Create(destPath)
	}

	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return nil, fmt.Errorf("cannot download file to %s since it's a directory", destPath)
	}

	return os.Create(destPath)
}
