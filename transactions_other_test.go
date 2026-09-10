package fio

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSetLastTransactionIDReturnsRequestError(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	err := client.SetLastTransactionID(context.Background(), 1)
	if err == nil || !strings.Contains(err.Error(), "500 Internal Server Error") {
		t.Fatalf("error = %v, want server error", err)
	}
}

func TestSetLastFailedTransactionDateReturnsRequestError(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	err := client.SetLastFailedTransactionDate(context.Background(), time.Now())
	if err == nil || !strings.Contains(err.Error(), "403 Forbidden") {
		t.Fatalf("error = %v, want forbidden error", err)
	}
}
