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

// GetAccountInfo returns metadata and balances for the account associated with
// the client's API token.
func (receiver *client) GetAccountInfo(ctx context.Context) (response.AccountInfo, error) {
	date := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	resp, err := receiver.request[response.TransactionsResponse](
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"/periods/{token}/%s/%s/transactions.json",
			date,
			date,
		),
		nil,
	)
	if err != nil {
		return response.AccountInfo{}, fmt.Errorf("failed issuing a request: %w", err)
	}

	return resp.AccountStatement.Info, nil
}
