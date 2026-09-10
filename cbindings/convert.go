package main

/*
#include <stdint.h>
#include <stdlib.h>
#include "headers/transactions.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/shopspring/decimal"
	"go.chrastecky.dev/fio-api/fio/dto"
	"go.chrastecky.dev/fio-api/fio/types"
)

func stringToC(value string) *C.char {
	return C.CString(value)
}

func optionalStringToC(value *string) *C.char {
	if value == nil {
		return nil
	}

	return C.CString(*value)
}

func optionalStringFromC(value *C.char) *string {
	if value == nil {
		return nil
	}
	result := C.GoString(value)

	return &result
}

func requiredStringFromC(name string, value *C.char) (string, error) {
	if value == nil {
		return "", nullPointerError(name)
	}

	return C.GoString(value), nil
}

func dateFromC(name string, value *C.char) (types.Date, error) {
	text, err := requiredStringFromC(name, value)
	if err != nil {
		return types.Date{}, err
	}

	var date types.Date
	if err = date.UnmarshalText([]byte(text)); err != nil {
		return types.Date{}, fmt.Errorf("invalid %s: %w", name, err)
	}

	return date, nil
}

func transactionToC(value dto.Transaction) C.FioTransaction {
	result := C.FioTransaction{
		id:                     C.int64_t(value.ID.Value),
		date:                   stringToC(value.Date.Value.String()),
		amount:                 stringToC(value.Amount.Value.String()),
		currency:               stringToC(value.Currency.Value),
		counterparty_account:   stringToC(value.CounterpartyAccount.Value),
		counterparty_name:      stringToC(value.CounterpartyName.Value),
		counterparty_bank_code: stringToC(value.CounterpartyBankCode.Value),
		counterparty_bank_name: stringToC(value.CounterpartyBankName.Value),
		constant_symbol:        optionalStringToC(value.ConstantSymbol.Value),
		variable_symbol:        optionalStringToC(value.VariableSymbol.Value),
		specific_symbol:        optionalStringToC(value.SpecificSymbol.Value),
		user_identity:          optionalStringToC(value.UserIdentity.Value),
		transaction_type:       stringToC(value.TransactionType.Value.String()),
		performed_by:           optionalStringToC(value.PerformedBy.Value),
		additional_info:        optionalStringToC(value.AdditionalInfo.Value),
		comment:                optionalStringToC(value.Comment.Value),
		bic:                    optionalStringToC(value.BIC.Value),
		payer_reference:        optionalStringToC(value.PayerReference.Value),
	}

	if value.InstructionID.Value != nil {
		result.instruction_id = (*C.int64_t)(C.malloc(C.size_t(unsafe.Sizeof(C.int64_t(0)))))
		if result.instruction_id != nil {
			*result.instruction_id = C.int64_t(*value.InstructionID.Value)
		}
	}

	return result
}

func freeTransaction(value *C.FioTransaction) {
	if value == nil {
		return
	}

	for _, ptr := range []*C.char{value.date, value.amount, value.currency, value.counterparty_account,
		value.counterparty_name, value.counterparty_bank_code, value.counterparty_bank_name,
		value.constant_symbol, value.variable_symbol, value.specific_symbol, value.user_identity,
		value.transaction_type, value.performed_by, value.additional_info, value.comment, value.bic,
		value.payer_reference} {
		C.free(unsafe.Pointer(ptr))
	}

	C.free(unsafe.Pointer(value.instruction_id))
	*value = C.FioTransaction{}
}

func domesticTransactionFromC(value C.FioDomesticTransaction) (dto.DomesticTransaction, error) {
	accountFrom, err := requiredStringFromC("transaction.account_from", value.account_from)
	if err != nil {
		return dto.DomesticTransaction{}, err
	}

	amountText, err := requiredStringFromC("transaction.amount", value.amount)
	if err != nil {
		return dto.DomesticTransaction{}, err
	}

	amount, err := decimal.NewFromString(amountText)
	if err != nil {
		return dto.DomesticTransaction{}, fmt.Errorf("invalid transaction.amount: %w", err)
	}

	accountTo, err := requiredStringFromC("transaction.account_to", value.account_to)
	if err != nil {
		return dto.DomesticTransaction{}, err
	}

	bankCode, err := requiredStringFromC("transaction.bank_code", value.bank_code)
	if err != nil {
		return dto.DomesticTransaction{}, err
	}

	var date types.Date
	if value.date != nil && C.GoString(value.date) != "" {
		date, err = dateFromC("transaction.date", value.date)
		if err != nil {
			return dto.DomesticTransaction{}, err
		}
	}

	result := dto.DomesticTransaction{
		AccountFrom:         accountFrom,
		Amount:              amount,
		AccountTo:           accountTo,
		BankCode:            bankCode,
		ConstantSymbol:      optionalStringFromC(value.constant_symbol),
		VariableSymbol:      optionalStringFromC(value.variable_symbol),
		SpecificSymbol:      optionalStringFromC(value.specific_symbol),
		MessageForRecipient: optionalStringFromC(value.message_for_recipient),
		Comment:             optionalStringFromC(value.comment),
		Date:                date,
	}

	if value.currency != nil {
		result.Currency = C.GoString(value.currency)
	}
	if value.payment_type != nil {
		result.PaymentType = dto.DomesticPaymentType(C.GoString(value.payment_type))
	}

	return result, nil
}
