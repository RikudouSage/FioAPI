package fio

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.chrastecky.dev/fio-api/fio/dto"
	"go.chrastecky.dev/fio-api/fio/internal/response"
)

func (receiver *client) TransactionsByDate(ctx context.Context, startDate time.Time, endDate time.Time) ([]dto.Transaction, error) {
	resp, err := receiver.request[response.TransactionsResponse](
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"/periods/{token}/%s/%s/transactions.json",
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed issuing a request: %w", err)
	}

	return resp.AccountStatement.TransactionList.Transactions, nil
}

func (receiver *client) TransactionsSinceLastPull(ctx context.Context) ([]dto.Transaction, error) {
	resp, err := receiver.request[response.TransactionsResponse](
		ctx,
		http.MethodGet,
		"/last/{token}/transactions.json",
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed issuing a request: %w", err)
	}

	return resp.AccountStatement.TransactionList.Transactions, nil
}

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
