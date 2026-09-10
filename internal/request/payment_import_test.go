package request

import (
	"testing"

	"go.chrastecky.dev/fio-api/fio/dto"
)

func TestPaymentImportProvideDefaults(t *testing.T) {
	transaction := &dto.DomesticTransaction{}
	request := PaymentImport{Orders: PaymentOrders{DomesticTransactions: []*dto.DomesticTransaction{transaction}}}
	request.ProvideDefaults()

	if request.XMLNSXSI != "http://www.w3.org/2001/XMLSchema-instance" {
		t.Errorf("XMLNSXSI = %q", request.XMLNSXSI)
	}
	if request.NoNamespaceSchemaLocation != "http://www.fio.cz/schema/importIB.xsd" {
		t.Errorf("NoNamespaceSchemaLocation = %q", request.NoNamespaceSchemaLocation)
	}
	if transaction.Currency != "CZK" || transaction.PaymentType != dto.DomesticPaymentStandard {
		t.Errorf("transaction defaults were not applied: %+v", transaction)
	}
}

func TestPaymentImportPreservesSchemaSettings(t *testing.T) {
	request := PaymentImport{XMLNSXSI: "custom namespace", NoNamespaceSchemaLocation: "custom schema"}
	request.ProvideDefaults()
	if request.XMLNSXSI != "custom namespace" || request.NoNamespaceSchemaLocation != "custom schema" {
		t.Errorf("schema settings were overwritten: %+v", request)
	}
}
