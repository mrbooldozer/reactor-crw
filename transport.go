package reactor_crw

import (
	"fmt"
	"io"
	"net/http"
)

// Headers defines a simple wrapper for client headers.
type Headers map[string]string

// HttpTransport allows making a network request over HTTP protocol. It is a
// wrapper over HTTP.Client with a set of default params required by the crawler.
type HttpTransport struct {
	headers Headers
	client  *http.Client
}

// NewHttpTransport creates a new *HttpTransport with provided client and custom
// headers. Custom headers can override the default ones.
func NewHttpTransport(c *http.Client, h Headers) *HttpTransport {
	var defaultHeaders = Headers{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:92.0) Gecko/20100101 Firefox/92.0",
		"Referer":    "http://joyreactor.cc/",
		"DNT":        "1",
	}

	t := &HttpTransport{
		headers: defaultHeaders,
		client: c,
	}

	for k, v := range h {
		t.headers[k] = v
	}

	return t
}

// FetchData makes an HTTP request using provided URL and returns the response
// body as io.ReadCloser interface. Each request will be prepared with provided
// headers list. The end client is responsible for closing the response body.
func (t *HttpTransport) FetchData(url string) (io.ReadCloser, error) {
	req, err := t.prepareRequest(http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	res, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot make request to %s: %w", url, err)
	}

	return res.Body, nil
}

func (t *HttpTransport) prepareRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request to %s: %w", url, err)
	}

	for k, v := range t.headers {
		req.Header.Add(k, v)
	}

	return req, nil
}
