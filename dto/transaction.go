package dto

import (
	"github.com/shopspring/decimal"
	. "go.chrastecky.dev/fio-api/fio/types"
)

// Transaction is one booked account transaction returned by Fio bank.
//
// Each field is wrapped in [TransactionValue] because the wire format
// includes column metadata alongside the actual value.
type Transaction struct {
	// ID is Fio bank's unique transaction identifier.
	ID TransactionValue[int64] `json:"column22"`
	// Date is the transaction's booking date and UTC offset.
	Date TransactionValue[TimezonedDate] `json:"column0"`
	// Amount is positive for credits and negative for debits.
	Amount TransactionValue[decimal.Decimal] `json:"column1"`
	// Currency is the ISO 4217 currency code of Amount.
	Currency TransactionValue[string] `json:"column14"`
	// CounterpartyAccount is the counterparty's account number.
	CounterpartyAccount TransactionValue[string] `json:"column2"`
	// CounterpartyName is the account holder name reported by the bank.
	CounterpartyName TransactionValue[string] `json:"column10"`
	// CounterpartyBankCode is the counterparty bank's domestic code.
	CounterpartyBankCode TransactionValue[string] `json:"column3"`
	// CounterpartyBankName is the counterparty bank's name.
	CounterpartyBankName TransactionValue[string] `json:"column12"`
	// ConstantSymbol is the optional constant symbol.
	ConstantSymbol TransactionValue[*string] `json:"column4"`
	// VariableSymbol is the optional variable symbol.
	VariableSymbol TransactionValue[*string] `json:"column5"`
	// SpecificSymbol is the optional specific symbol.
	SpecificSymbol TransactionValue[*string] `json:"column6"`
	// UserIdentity is the optional user-supplied transaction identifier.
	UserIdentity TransactionValue[*string] `json:"column7"`
	// TransactionType classifies the transaction.
	TransactionType TransactionValue[TransactionType] `json:"column8"`
	// PerformedBy identifies who performed the transaction, when supplied.
	PerformedBy TransactionValue[*string] `json:"column9"`
	// AdditionalInfo contains optional information supplied by the bank.
	AdditionalInfo TransactionValue[*string] `json:"column18"`
	// Comment is the optional transaction comment.
	Comment TransactionValue[*string] `json:"column25"`
	// BIC is the optional counterparty bank identifier code.
	BIC TransactionValue[*string] `json:"column26"`
	// InstructionID identifies the payment instruction, when available.
	InstructionID TransactionValue[*int64] `json:"column17"`
	// PayerReference is the optional payer-provided reference.
	PayerReference TransactionValue[*string] `json:"column27"`
}
