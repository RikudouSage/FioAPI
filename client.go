package fio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.chrastecky.dev/fio-api/fio/dto"
	. "go.chrastecky.dev/fio-api/fio/internal/helper"
)

type Client interface {
	TransactionsByDate(ctx context.Context, startDate time.Time, endDate time.Time) ([]dto.Transaction, error)
	TransactionsSinceLastPull(ctx context.Context) ([]dto.Transaction, error)

	SetLastTransactionID(ctx context.Context, id int64) error
	SetLastFailedTransactionDate(ctx context.Context, date time.Time) error
}

type client struct {
	baseURL string
	token   string
	debug   bool
}

func NewClient(token string, options ...Option) (Client, error) {
	instance := &client{
		token: token,
	}

	for _, option := range options {
		if err := option(instance); err != nil {
			return nil, fmt.Errorf("failed applying option: %w", err)
		}
	}

	if instance.baseURL == "" {
		instance.baseURL = "https://fioapi.fio.cz/v1/rest"
	}

	return instance, nil
}

func (receiver *client) request[TResult any](
	ctx context.Context,
	method string,
	path string,
	body any,
) (TResult, error) {
	var result TResult

	var bodyReader io.Reader
	if body != nil {
		switch body := body.(type) {
		case io.Reader:
			bodyReader = body
		case []byte:
			bodyReader = bytes.NewReader(body)
		case string:
			bodyReader = strings.NewReader(body)
		default:
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				if err != nil {
					return result, fmt.Errorf("failed marshalling body as JSON: %w", err)
				}
			}
			bodyReader = bytes.NewReader(bodyBytes)
		}
	}

	uri := receiver.uriWithPath(path)
	uri = strings.ReplaceAll(uri, "%7Btoken%7D", receiver.token)

	req, err := http.NewRequestWithContext(ctx, method, uri, bodyReader)
	if err != nil {
		return result, fmt.Errorf("failed creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed sending request: %w", err)
	}
	defer DrainBody(resp.Body)

	if resp.StatusCode < minSuccessfulStatusCode || resp.StatusCode > maxSuccessfulStatusCode {
		if resp.StatusCode == http.StatusTooManyRequests {
			return result, ErrTooSoon
		}

		errMsg := fmt.Sprintf("failed sending request: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		if receiver.debug {
			rawResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				rawResponse = []byte(fmt.Sprintf("failed reading response body: %v", err))
			}
			errMsg += fmt.Sprintf("\nBody: %s", string(rawResponse))
		}

		return result, errors.New(errMsg)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil && !errors.Is(err, io.EOF) {
		return result, fmt.Errorf("failed parsing response: %w", err)
	}

	return result, nil
}

func (receiver *client) uriWithPath(path string) string {
	uri := lo.Must(url.Parse(receiver.baseURL))
	if uri.Path != "" {
		uri.Path = strings.Join([]string{
			strings.TrimSuffix(uri.Path, "/"),
			strings.TrimPrefix(path, "/"),
		}, "/")
	}

	return uri.String()
}
