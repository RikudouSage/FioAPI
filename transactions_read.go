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
func (receiver *client) GetAccountInfo(ctx context.Context) (dto.AccountInfo, error) {
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
		return dto.AccountInfo{}, fmt.Errorf("failed issuing a request: %w", err)
	}

	return resp.AccountStatement.Info, nil
}

// AccountInfoAndTransactions returns account metadata and transactions added
// since the token's last successful download.
func (receiver *client) AccountInfoAndTransactions(ctx context.Context) (dto.AccountInfo, []dto.Transaction, error) {
	resp, err := receiver.request[response.TransactionsResponse](
		ctx,
		http.MethodGet,
		"/last/{token}/transactions.json",
		nil,
	)
	if err != nil {
		return dto.AccountInfo{}, nil, fmt.Errorf("failed issuing a request: %w", err)
	}

	return resp.AccountStatement.Info, resp.AccountStatement.TransactionList.Transactions, nil
}

// AccountInfoAndTransactionsByDate returns account metadata and transactions
// booked in the inclusive date range from startDate through endDate.
func (receiver *client) AccountInfoAndTransactionsByDate(ctx context.Context, startDate time.Time, endDate time.Time) (dto.AccountInfo, []dto.Transaction, error) {
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
		return dto.AccountInfo{}, nil, fmt.Errorf("failed issuing a request: %w", err)
	}

	return resp.AccountStatement.Info, resp.AccountStatement.TransactionList.Transactions, nil
}
