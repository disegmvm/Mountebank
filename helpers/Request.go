package helpers

import (
	"io"
	"net/http"
)

type Request struct {
	http *http.Request
	path string
	Err  error
}

const (
	correlationIDHeader = "x-correlation-id"
	testingIDHeader     = "x-testing-id"
)

func NewRequest(baseURL string, method string, path string) *Request {
	h, err := http.NewRequest(method, baseURL+path, nil)

	req := Request{
		http: nil,
		path: path,
		Err:  nil,
	}

	if err != nil {
		req.Err = err

		return &req
	}

	req.http = h

	return &req
}

func (r *Request) WithHeaders(headers map[string]string) *Request {
	if r.Err != nil {
		return r
	}
	for k, v := range headers {
		r.http.Header.Add(k, v)
	}

	return r
}

func (r *Request) WithBody(body io.Reader) *Request {
	if r.Err != nil {
		return r
	}
	r.http.Body = io.NopCloser(body)

	return r
}

func (r *Request) WithContentTypeHeaderJSON() *Request {
	if r.Err != nil {
		return r
	}
	r.http.Header.Add("Content-Type", "application/json")

	return r
}

func (r *Request) WithCorrelationID(cid string) *Request {
	if r.Err != nil {
		return r
	}
	r.http.Header.Add(correlationIDHeader, cid)

	return r
}

func (r *Request) WithTestingID(tid string) *Request {
	if r.Err != nil {
		return r
	}
	r.http.Header.Add(testingIDHeader, tid)

	return r
}
