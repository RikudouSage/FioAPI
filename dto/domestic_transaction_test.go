package dto

import (
	"testing"
	"time"

	"go.chrastecky.dev/fio-api/fio/types"
)

func TestDomesticTransactionProvideDefaults(t *testing.T) {
	before := time.Now().Add(-time.Second)
	var transaction DomesticTransaction
	transaction.ProvideDefaults()
	after := time.Now().Add(time.Second)

	if transaction.Currency != "CZK" {
		t.Errorf("Currency = %q, want CZK", transaction.Currency)
	}
	if transaction.PaymentType != DomesticPaymentStandard {
		t.Errorf("PaymentType = %q, want %q", transaction.PaymentType, DomesticPaymentStandard)
	}
	date := transaction.Date.AsTime()
	if date.Before(before) || date.After(after) {
		t.Errorf("Date = %v, want current time", date)
	}
}

func TestDomesticTransactionDefaultsDoNotOverwriteValues(t *testing.T) {
	wantDate := types.Date(time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC))
	transaction := DomesticTransaction{
		Currency:    "EUR",
		Date:        wantDate,
		PaymentType: DomesticPaymentPriority,
	}
	transaction.ProvideDefaults()
	if transaction.Currency != "EUR" || transaction.Date != wantDate || transaction.PaymentType != DomesticPaymentPriority {
		t.Errorf("ProvideDefaults() overwrote values: %+v", transaction)
	}
}
