package fio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc, options ...Option) (Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	options = append([]Option{WithBaseURL(server.URL)}, options...)
	client, err := NewClient("test-token", options...)
	if err != nil {
		server.Close()
		t.Fatalf("NewClient() error = %v", err)
	}
	t.Cleanup(server.Close)
	return client, server
}

func TestNewClientRejectsInvalidBaseURL(t *testing.T) {
	_, err := NewClient("token", WithBaseURL("localhost:8080"))
	if err == nil {
		t.Fatal("NewClient() unexpectedly accepted a URL without scheme and host")
	}
}

func TestTransactionsSinceLastPullUsesBaseURLWithoutPath(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/last/test-token/transactions.json" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/last/test-token/transactions.json")
		}
		_, _ = io.WriteString(w, `{"accountStatement":{"transactionList":{"transaction":[]}}}`)
	})

	transactions, err := client.TransactionsSinceLastPull(context.Background())
	if err != nil {
		t.Fatalf("TransactionsSinceLastPull() error = %v", err)
	}
	if len(transactions) != 0 {
		t.Fatalf("got %d transactions, want 0", len(transactions))
	}
}

func TestTransactionsByDateBuildsPathAndDecodesResponse(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		want := "/periods/test-token/2025-01-02/2025-01-31/transactions.json"
		if r.URL.Path != want {
			t.Errorf("path = %q, want %q", r.URL.Path, want)
		}
		_, _ = io.WriteString(w, `{"accountStatement":{"transactionList":{"transaction":[{`+
			`"column22":{"value":42},"column0":{"value":"2025-01-03+0100"},`+
			`"column1":{"value":"12.50"},"column14":{"value":"CZK"},`+
			`"column8":{"value":"Okamžitá příchozí platba"}}]}}}`)
	})

	start := time.Date(2025, 1, 2, 23, 0, 0, 0, time.FixedZone("test", -5*60*60))
	end := time.Date(2025, 1, 31, 1, 0, 0, 0, time.UTC)
	transactions, err := client.TransactionsByDate(context.Background(), start, end)
	if err != nil {
		t.Fatalf("TransactionsByDate() error = %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("got %d transactions, want 1", len(transactions))
	}
	got := transactions[0]
	if got.ID.Value != 42 || got.Amount.Value.String() != "12.5" || got.Currency.Value != "CZK" {
		t.Errorf("decoded transaction = ID %d, amount %s, currency %q", got.ID.Value, got.Amount.Value, got.Currency.Value)
	}
	if !got.TransactionType.Value.IsIncoming() || !got.TransactionType.Value.IsImmediate() {
		t.Errorf("transaction type %q was not classified correctly", got.TransactionType.Value)
	}
}

func TestSetLastMarkers(t *testing.T) {
	wants := []string{
		"/set-last-id/test-token/123456789/",
		"/set-last-date/test-token/2025-02-03/",
	}
	request := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if request >= len(wants) {
			t.Fatalf("unexpected extra request to %q", r.URL.Path)
		}
		if r.URL.Path != wants[request] {
			t.Errorf("request %d path = %q, want %q", request, r.URL.Path, wants[request])
		}
		request++
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.SetLastTransactionID(context.Background(), 123456789); err != nil {
		t.Fatalf("SetLastTransactionID() error = %v", err)
	}
	date := time.Date(2025, 2, 3, 22, 10, 0, 0, time.FixedZone("test", 9*60*60))
	if err := client.SetLastFailedTransactionDate(context.Background(), date); err != nil {
		t.Fatalf("SetLastFailedTransactionDate() error = %v", err)
	}
}

func TestHTTPErrorHandling(t *testing.T) {
	t.Run("conflict", func(t *testing.T) {
		client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusConflict)
		})
		_, err := client.TransactionsSinceLastPull(context.Background())
		if !errors.Is(err, ErrTooSoon) {
			t.Fatalf("error = %v, want ErrTooSoon", err)
		}
	})

	t.Run("debug body", func(t *testing.T) {
		client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, "bank diagnostic")
		}, WithDebug(true))
		_, err := client.TransactionsSinceLastPull(context.Background())
		if err == nil || !strings.Contains(err.Error(), "400 Bad Request") || !strings.Contains(err.Error(), "bank diagnostic") {
			t.Fatalf("error = %v, want status and response body", err)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			_, _ = io.WriteString(w, "not JSON")
		})
		_, err := client.TransactionsSinceLastPull(context.Background())
		if err == nil || !strings.Contains(err.Error(), "failed parsing response") {
			t.Fatalf("error = %v, want JSON parsing error", err)
		}
	})
}

func TestRequestHonorsCanceledContext(t *testing.T) {
	client, _ := newTestClient(t, func(http.ResponseWriter, *http.Request) {
		t.Error("server handler should not be reached for an already canceled context")
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.TransactionsSinceLastPull(ctx)
	if err == nil || !strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Fatalf("error = %v, want context cancellation", err)
	}
}

func TestBaseURLPathIsPreserved(t *testing.T) {
	var got string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Path
		_, _ = fmt.Fprint(w, `{}`)
	}))
	defer server.Close()
	client, err := NewClient("test-token", WithBaseURL(server.URL+"/api/v1/"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.TransactionsSinceLastPull(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := "/api/v1/last/test-token/transactions.json"; got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
}
