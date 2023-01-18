package helpers

import (
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type HttpHelper struct {
	BaseURL string
	Client  Client
}

func NewHttpHelper(baseURL string) *HttpHelper {
	return &HttpHelper{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (h *HttpHelper) POST(path string) *Request {
	return h.Request("POST", path)
}

func (h *HttpHelper) GET(path string) *Request {
	return h.Request("GET", path)
}

func (h *HttpHelper) Request(method string, path string) *Request {
	req := NewRequest(h.BaseURL, method, path)

	return req
}

func (h *HttpHelper) Send(r *Request) (*http.Response, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return h.Client.Do(r.http)
}
