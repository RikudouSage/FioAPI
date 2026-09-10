package fio

import (
	"context"
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
