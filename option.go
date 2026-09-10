package fio

import (
	"fmt"
	"net/url"
	"strings"
)

// Option configures a Client created by [NewClient].
type Option func(instance *client) error

// WithBaseURL sets the base URL used for API requests.
// The URL must contain a scheme and host. This option is primarily useful for
// proxies, compatible API implementations, and tests.
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

// WithDebug enables or disables debug error details and formatted XML request
// bodies. Debug errors may contain response bodies supplied by the server.
func WithDebug(debug bool) Option {
	return func(instance *client) error {
		instance.debug = debug
		return nil
	}
}
