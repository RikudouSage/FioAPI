package fio

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (receiver *client) SetLastTransactionID(ctx context.Context, id int64) error {
	_, err := receiver.request[any](
		ctx,
		http.MethodGet,
		fmt.Sprintf("/set-last-id/{token}/%d/", id),
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed issuing a request: %w", err)
	}

	return nil
}

func (receiver *client) SetLastFailedTransactionDate(ctx context.Context, date time.Time) error {
	_, err := receiver.request[any](
		ctx,
		http.MethodGet,
		fmt.Sprintf("/set-last-date/{token}/%s/", date.Format("2006-01-02")),
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed issuing a request: %w", err)
	}

	return nil
}
