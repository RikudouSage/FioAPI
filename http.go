package fio

import "errors"

var ErrTooSoon = errors.New("the request was issued too soon")

const (
	minSuccessfulStatusCode = 200
	maxSuccessfulStatusCode = 299
)
