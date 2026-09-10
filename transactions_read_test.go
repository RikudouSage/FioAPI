package fio

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestTransactionsByDateReturnsRequestError(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	})
	_, err := client.TransactionsByDate(context.Background(), time.Now(), time.Now())
	if err == nil || !strings.Contains(err.Error(), "503 Service Unavailable") {
		t.Fatalf("error = %v, want service-unavailable error", err)
	}
}

func TestTransactionsResponseWithoutListReturnsNilSlice(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"accountStatement":{}}`))
	})
	transactions, err := client.TransactionsSinceLastPull(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if transactions != nil {
		t.Errorf("transactions = %#v, want nil", transactions)
	}
}

func TestAccountInfoAndTransactionsReturnsAccountInfoAndTransactions(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if want := "/last/test-token/transactions.json"; r.URL.Path != want {
			t.Errorf("path = %q, want %q", r.URL.Path, want)
		}
		_, _ = io.WriteString(w, `{"accountStatement":{`+
			`"info":{"accountId":"1234567890","bankId":"2010","currency":"CZK",`+
			`"iban":"CZ0120100000001234567890","bic":"FIOBCZPPXXX",`+
			`"openingBalance":"100.25","closingBalance":"112.75"},`+
			`"transactionList":{"transaction":[{`+
			`"column22":{"value":42},"column1":{"value":"12.50"},`+
			`"column14":{"value":"CZK"}}]}}}`)
	})

	info, transactions, err := client.AccountInfoAndTransactions(context.Background())
	if err != nil {
		t.Fatalf("AccountInfoAndTransactions() error = %v", err)
	}
	if info.AccountID != "1234567890" || info.BankID != "2010" || info.Currency != "CZK" {
		t.Errorf("account info = %#v, want account 1234567890/2010 in CZK", info)
	}
	if info.OpeningBalance.String() != "100.25" || info.ClosingBalance.String() != "112.75" {
		t.Errorf("balances = %s/%s, want 100.25/112.75", info.OpeningBalance, info.ClosingBalance)
	}
	if len(transactions) != 1 {
		t.Fatalf("got %d transactions, want 1", len(transactions))
	}
	if got := transactions[0]; got.ID.Value != 42 || got.Amount.Value.String() != "12.5" || got.Currency.Value != "CZK" {
		t.Errorf("transaction = %#v, want ID 42, amount 12.5 CZK", got)
	}
}

func TestAccountInfoAndTransactionsReturnsRequestError(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	})

	info, transactions, err := client.AccountInfoAndTransactions(context.Background())
	if err == nil || !strings.Contains(err.Error(), "503 Service Unavailable") {
		t.Fatalf("error = %v, want service-unavailable error", err)
	}
	if info.AccountID != "" {
		t.Errorf("account info = %#v, want zero value", info)
	}
	if transactions != nil {
		t.Errorf("transactions = %#v, want nil", transactions)
	}
}
