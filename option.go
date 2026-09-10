package fio

import (
	"fmt"
	"net/url"
	"strings"
)

type Option func(instance *client) error

func WithBaseURL(uri string) Option {
	return func(instance *client) error {
		parsed, err := url.Parse(uri)
		if err != nil {
			return fmt.Errorf("failed parsing base url: %w", err)
		}

		if parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("failed parsing base url '%s': missing scheme, or host", uri)
		}

		instance.baseURL = strings.TrimSuffix(uri, "/")
		return nil
	}
}

func WithDebug(debug bool) Option {
	return func(instance *client) error {
		instance.debug = debug
		return nil
	}
}
