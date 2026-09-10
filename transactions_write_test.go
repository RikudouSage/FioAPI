package fio

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"go.chrastecky.dev/fio-api/fio/dto"
)

func TestIssueDomesticTransactionCreatesMultipartXML(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/import/" {
			t.Errorf("request = %s %s, want POST /import/", r.Method, r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm() error = %v", err)
		}
		if got := r.FormValue("type"); got != "xml" {
			t.Errorf("type = %q, want xml", got)
		}
		if got := r.FormValue("token"); got != "test-token" {
			t.Errorf("token = %q, want test-token", got)
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile() error = %v", err)
		}
		defer file.Close()
		if header.Filename != "payments.xml" {
			t.Errorf("filename = %q, want payments.xml", header.Filename)
		}
		data, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		var order struct {
			XMLNS string `xml:"xsi,attr"`
			Items []struct {
				AccountFrom string `xml:"accountFrom"`
				Currency    string `xml:"currency"`
				Amount      string `xml:"amount"`
				AccountTo   string `xml:"accountTo"`
				BankCode    string `xml:"bankCode"`
				PaymentType string `xml:"paymentType"`
			} `xml:"Orders>DomesticTransaction"`
		}
		if err := xml.Unmarshal(data, &order); err != nil {
			t.Fatalf("generated XML is invalid: %v\n%s", err, data)
		}
		if len(order.Items) != 1 {
			t.Fatalf("got %d payment orders, want 1\n%s", len(order.Items), data)
		}
		got := order.Items[0]
		if got.AccountFrom != "123" || got.AccountTo != "456" || got.BankCode != "2010" || got.Amount != "99.95" {
			t.Errorf("unexpected payment XML: %+v", got)
		}
		if got.Currency != "CZK" || got.PaymentType != string(dto.DomesticPaymentStandard) {
			t.Errorf("defaults not applied: currency=%q paymentType=%q", got.Currency, got.PaymentType)
		}
		if !strings.HasPrefix(string(data), xml.Header) {
			t.Error("payment XML does not start with the XML declaration")
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = io.WriteString(w, `<responseImport><result><errorCode>0</errorCode><status>ok</status></result></responseImport>`)
	})

	err := client.IssueDomesticTransaction(context.Background(), dto.DomesticTransaction{
		AccountFrom: "123",
		Amount:      decimal.RequireFromString("99.95"),
		AccountTo:   "456",
		BankCode:    "2010",
	})
	if err != nil {
		t.Fatalf("IssueDomesticTransaction() error = %v", err)
	}
}

func TestIssueDomesticTransactionReportsImportFailure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `<responseImport><result><errorCode>11</errorCode><status>error</status></result></responseImport>`)
	})
	err := client.IssueDomesticTransaction(context.Background(), dto.DomesticTransaction{})
	if err == nil || !strings.Contains(err.Error(), "status error (code: 11)") {
		t.Fatalf("error = %v, want import status and code", err)
	}
}
