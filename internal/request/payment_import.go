package request

import (
	"encoding/xml"

	"go.chrastecky.dev/fio-api/fio/dto"
)

type PaymentImport struct {
	XMLName                   xml.Name `xml:"Import"`
	XMLNSXSI                  string   `xml:"xmlns:xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"xsi:noNamespaceSchemaLocation,attr"`

	Orders PaymentOrders `xml:"Orders"`
}

type PaymentOrders struct {
	DomesticTransactions []*dto.DomesticTransaction `xml:"DomesticTransaction"`
}

func (receiver *PaymentImport) ProvideDefaults() {
	if receiver.XMLNSXSI == "" {
		receiver.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"
	}
	if receiver.NoNamespaceSchemaLocation == "" {
		receiver.NoNamespaceSchemaLocation = "http://www.fio.cz/schema/importIB.xsd"
	}

	for _, transaction := range receiver.Orders.DomesticTransactions {
		transaction.ProvideDefaults()
	}
}
