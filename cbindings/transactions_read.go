package main

/*
#include <stdlib.h>
#include "headers/transactions.h"
*/
import "C"

import (
	"unsafe"

	"go.chrastecky.dev/fio-api/fio/dto"
)

func accountInfoToC(value dto.AccountInfo) C.FioAccountInfo {
	result := C.FioAccountInfo{
		account_id:      stringToC(value.AccountID),
		bank_id:         stringToC(value.BankID),
		currency:        stringToC(value.Currency),
		iban:            stringToC(value.IBAN),
		bic:             stringToC(value.BIC),
		opening_balance: stringToC(value.OpeningBalance.String()),
		closing_balance: stringToC(value.ClosingBalance.String()),
		date_start:      stringToC(value.DateStart.String()),
		date_end:        stringToC(value.DateEnd.String()),
	}
	if value.YearList != nil {
		result.year_list = (*C.uint16_t)(C.malloc(C.size_t(unsafe.Sizeof(C.uint16_t(0)))))
		if result.year_list != nil {
			*result.year_list = C.uint16_t(*value.YearList)
		}
	}
	setOptionalInt64 := func(out **C.int64_t, value *int64) {
		if value == nil {
			return
		}
		*out = (*C.int64_t)(C.malloc(C.size_t(unsafe.Sizeof(C.int64_t(0)))))
		if *out != nil {
			**out = C.int64_t(*value)
		}
	}
	setOptionalInt64(&result.id_list, value.IDList)
	setOptionalInt64(&result.id_from, value.IDFrom)
	setOptionalInt64(&result.id_to, value.IDTo)
	setOptionalInt64(&result.id_last_download, value.IDLastDownload)
	return result
}

func returnTransactions(out *C.FioTransactions, transactionsGoLen int, fill func(int) C.FioTransaction) C.FioResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioFailure
	}

	*out = C.FioTransactions{}
	if transactionsGoLen == 0 {
		clearLastError()
		return C.FioSuccess
	}

	size := unsafe.Sizeof(C.FioTransaction{})
	items := C.calloc(C.size_t(transactionsGoLen), C.size_t(size))
	if items == nil {
		setLastErrorMessage("failed allocating transactions")
		return C.FioFailure
	}

	slice := unsafe.Slice((*C.FioTransaction)(items), transactionsGoLen)
	for i := range slice {
		slice[i] = fill(i)
	}
	out.items = (*C.FioTransaction)(items)
	out.length = C.size_t(transactionsGoLen)

	clearLastError()
	return C.FioSuccess
}

// FioTransactionsByDate returns transactions in the inclusive date range.
// Dates use YYYY-MM-DD format. The caller owns
// the returned allocation and must release it with FioFreeTransactions.
//
//export FioTransactionsByDate
func FioTransactionsByDate(client C.ClientHandle, ctx C.ContextHandle, startDate, endDate *C.char, out *C.FioTransactions) C.FioResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioFailure
	}

	*out = C.FioTransactions{}
	clientGo, ctxGo, err := getCommonHandles(client, ctx)

	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	startDateGo, err := dateFromC("start_date", startDate)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}
	endDateGo, err := dateFromC("end_date", endDate)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	transactions, err := clientGo.TransactionsByDate(ctxGo, startDateGo.AsTime(), endDateGo.AsTime())
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	return returnTransactions(out, len(transactions), func(i int) C.FioTransaction {
		return transactionToC(transactions[i])
	})
}

// FioTransactionsSinceLastPull returns transactions added since the client's
// last successful download. The caller owns the returned allocation and must
// release it with FioFreeTransactions.
//
//export FioTransactionsSinceLastPull
func FioTransactionsSinceLastPull(client C.ClientHandle, ctx C.ContextHandle, out *C.FioTransactions) C.FioResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioFailure
	}
	*out = C.FioTransactions{}

	clientGo, ctxGo, err := getCommonHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	transactions, err := clientGo.TransactionsSinceLastPull(ctxGo)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	return returnTransactions(out, len(transactions), func(i int) C.FioTransaction {
		return transactionToC(transactions[i])
	})
}

// FioGetAccountInfo returns metadata and balances for the client's account.
// The caller owns the returned value and must release it with
// FioFreeAccountInfo.
//
//export FioGetAccountInfo
func FioGetAccountInfo(client C.ClientHandle, ctx C.ContextHandle, out *C.FioAccountInfo) C.FioResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioFailure
	}
	*out = C.FioAccountInfo{}

	clientGo, ctxGo, err := getCommonHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}
	info, err := clientGo.GetAccountInfo(ctxGo)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	*out = accountInfoToC(info)
	clearLastError()
	return C.FioSuccess
}

// FioFreeAccountInfo releases all memory owned by accountInfo. Passing NULL is
// safe.
//
//export FioFreeAccountInfo
func FioFreeAccountInfo(accountInfo *C.FioAccountInfo) {
	if accountInfo == nil {
		return
	}
	for _, ptr := range []*C.char{accountInfo.account_id, accountInfo.bank_id, accountInfo.currency,
		accountInfo.iban, accountInfo.bic, accountInfo.opening_balance, accountInfo.closing_balance,
		accountInfo.date_start, accountInfo.date_end} {
		C.free(unsafe.Pointer(ptr))
	}
	for _, ptr := range []unsafe.Pointer{unsafe.Pointer(accountInfo.year_list), unsafe.Pointer(accountInfo.id_list),
		unsafe.Pointer(accountInfo.id_from), unsafe.Pointer(accountInfo.id_to), unsafe.Pointer(accountInfo.id_last_download)} {
		C.free(ptr)
	}
	*accountInfo = C.FioAccountInfo{}
}

// FioFreeTransactions releases a transaction collection and all strings and
// optional values owned by it. Passing NULL is safe.
//
//export FioFreeTransactions
func FioFreeTransactions(transactions *C.FioTransactions) {
	if transactions == nil {
		return
	}

	if transactions.items != nil {
		items := unsafe.Slice(transactions.items, int(transactions.length))
		for i := range items {
			freeTransaction(&items[i])
		}
		C.free(unsafe.Pointer(transactions.items))
	}

	*transactions = C.FioTransactions{}
}
