package dto

import (
	"github.com/shopspring/decimal"
	. "go.chrastecky.dev/fio-api/fio/types"
)

type Transaction struct {
	ID                   TransactionValue[int64]           `json:"column22"`
	Date                 TransactionValue[TimezonedDate]   `json:"column0"`
	Amount               TransactionValue[decimal.Decimal] `json:"column1"`
	Currency             TransactionValue[string]          `json:"column14"`
	CounterpartyAccount  TransactionValue[string]          `json:"column2"`
	CounterpartyName     TransactionValue[string]          `json:"column10"`
	CounterpartyBankCode TransactionValue[string]          `json:"column3"`
	CounterpartyBankName TransactionValue[string]          `json:"column12"`
	ConstantSymbol       TransactionValue[*string]         `json:"column4"`
	VariableSymbol       TransactionValue[*string]         `json:"column5"`
	SpecificSymbol       TransactionValue[*string]         `json:"column6"`
	UserIdentity         TransactionValue[*string]         `json:"column7"`
	TransactionType      TransactionValue[TransactionType] `json:"column8"`
	PerformedBy          TransactionValue[*string]         `json:"column9"`
	AdditionalInfo       TransactionValue[*string]         `json:"column18"`
	Comment              TransactionValue[*string]         `json:"column25"`
	BIC                  TransactionValue[*string]         `json:"column26"`
	InstructionID        TransactionValue[*int64]          `json:"column17"`
	PayerReference       TransactionValue[*string]         `json:"column27"`
}
