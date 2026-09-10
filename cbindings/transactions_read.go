package main

/*
#include <stdlib.h>
#include "headers/transactions.h"
*/
import "C"

import "unsafe"

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
// Dates are expressed as milliseconds since the Unix epoch. The caller owns
// the returned allocation and must release it with FioFreeTransactions.
//
//export FioTransactionsByDate
func FioTransactionsByDate(client C.ClientHandle, ctx C.ContextHandle, startDateMs, endDateMs C.uint64_t, out *C.FioTransactions) C.FioResult {
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

	transactions, err := clientGo.TransactionsByDate(ctxGo, dateFromC(startDateMs), dateFromC(endDateMs))
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
