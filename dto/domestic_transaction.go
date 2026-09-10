package dto

import (
	"time"

	"github.com/shopspring/decimal"
	. "go.chrastecky.dev/fio-api/fio/types"
)

// DomesticPaymentType identifies the processing mode of a domestic payment.
type DomesticPaymentType string

const (
	// DomesticPaymentStandard requests standard processing. Immediate/instant payments should also use this, the bank
	// automatically marks them as immediate if it can.
	DomesticPaymentStandard DomesticPaymentType = "431001"
	// DomesticPaymentPriority requests priority processing.
	DomesticPaymentPriority DomesticPaymentType = "431005"
	// DomesticPaymentDirectDebit requests a direct debit.
	DomesticPaymentDirectDebit DomesticPaymentType = "431022"
)

// DomesticTransaction describes a Czech domestic payment order.
// Pointer fields are optional; nil omits the corresponding XML element.
type DomesticTransaction struct {
	// AccountFrom is the account number from which funds are sent.
	AccountFrom string `xml:"accountFrom"`
	// Currency is the ISO 4217 currency code. It defaults to CZK.
	Currency string `xml:"currency"`
	// Amount is the amount to transfer.
	Amount decimal.Decimal `xml:"amount"`
	// AccountTo is the recipient's account number.
	AccountTo string `xml:"accountTo"`
	// BankCode is the recipient bank's four-digit Czech bank code.
	BankCode string `xml:"bankCode"`
	// ConstantSymbol is the optional constant symbol.
	ConstantSymbol *string `xml:"ks,omitempty"`
	// VariableSymbol is the optional variable symbol.
	VariableSymbol *string `xml:"vs,omitempty"`
	// SpecificSymbol is the optional specific symbol.
	SpecificSymbol *string `xml:"ss,omitempty"`
	// Date is the requested payment date. It defaults to the current date.
	Date Date `xml:"date"`
	// MessageForRecipient is an optional message shown to the recipient.
	MessageForRecipient *string `xml:"messageForRecipient,omitempty"`
	// Comment is an optional payer-side comment.
	Comment *string `xml:"comment,omitempty"`
	// PaymentType controls payment processing. It defaults to
	// DomesticPaymentStandard.
	PaymentType DomesticPaymentType `xml:"paymentType"`
}

// ProvideDefaults fills empty Currency, Date, and PaymentType fields with CZK,
// the current date, and [DomesticPaymentStandard], respectively.
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
