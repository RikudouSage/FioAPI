package dto

import (
	"time"

	"github.com/shopspring/decimal"
	. "go.chrastecky.dev/fio-api/fio/types"
)

type DomesticPaymentType string

const (
	DomesticPaymentStandard    DomesticPaymentType = "431001"
	DomesticPaymentPriority    DomesticPaymentType = "431005"
	DomesticPaymentDirectDebit DomesticPaymentType = "431022"
)

type DomesticTransaction struct {
	AccountFrom         string              `xml:"accountFrom"`
	Currency            string              `xml:"currency"`
	Amount              decimal.Decimal     `xml:"amount"`
	AccountTo           string              `xml:"accountTo"`
	BankCode            string              `xml:"bankCode"`
	ConstantSymbol      *string             `xml:"ks,omitempty"`
	VariableSymbol      *string             `xml:"vs,omitempty"`
	SpecificSymbol      *string             `xml:"ss,omitempty"`
	Date                Date                `xml:"date"`
	MessageForRecipient *string             `xml:"messageForRecipient,omitempty"`
	Comment             *string             `xml:"comment,omitempty"`
	PaymentType         DomesticPaymentType `xml:"paymentType"`
}

func (receiver *DomesticTransaction) ProvideDefaults() {
	if receiver.Currency == "" {
		receiver.Currency = "CZK"
	}
	if receiver.Date == Date(time.Time{}) {
		receiver.Date = Date(time.Now())
	}
	if receiver.PaymentType == "" {
		receiver.PaymentType = DomesticPaymentStandard
	}
}
