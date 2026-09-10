package fio

import (
	"errors"
	"io"
	"net/http"
)

// ErrTooSoon is returned when Fio bank rejects a request because the token's
// request frequency limit has been exceeded.
var ErrTooSoon = errors.New("the request was issued too soon")
var ErrUnauthorizedLongAccess = errors.New("the api token is not authorized to get data older than 90 days")

const (
	minSuccessfulStatusCode = 200
	maxSuccessfulStatusCode = 299
)

type responseBodyDecoder func(body io.Reader) (any, error)
type httpOptions struct {
	responseDecoder responseBodyDecoder
}

type httpOption func(req *http.Request, options *httpOptions) error

func withHeader(name, value string) httpOption {
	return func(req *http.Request, _ *httpOptions) error {
		req.Header.Set(name, value)
		return nil
	}
}

func withResponseBodyDecoder(parser responseBodyDecoder) httpOption {
	return func(_ *http.Request, options *httpOptions) error {
		options.responseDecoder = parser
		return nil
	}
}
