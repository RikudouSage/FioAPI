package main

/*
#include <stdint.h>
#include "headers/common.h"
*/
import "C"

// FioSetLastTransactionID sets the transaction from which the client's next
// last-pull request continues.
//
//export FioSetLastTransactionID
func FioSetLastTransactionID(client C.ClientHandle, ctx C.ContextHandle, id C.int64_t) C.FioResult {
	clientGo, ctxGo, err := getCommonHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	if err = clientGo.SetLastTransactionID(ctxGo, int64(id)); err != nil {
		setLastError(err)
		return C.FioFailure
	}

	clearLastError()
	return C.FioSuccess
}

// FioSetLastFailedTransactionDate moves the client's last-pull marker to the
// supplied date in YYYY-MM-DD format.
//
//export FioSetLastFailedTransactionDate
func FioSetLastFailedTransactionDate(client C.ClientHandle, ctx C.ContextHandle, date *C.char) C.FioResult {
	clientGo, ctxGo, err := getCommonHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	dateGo, err := dateFromC("date", date)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	if err = clientGo.SetLastFailedTransactionDate(ctxGo, dateGo.AsTime()); err != nil {
		setLastError(err)
		return C.FioFailure
	}

	clearLastError()
	return C.FioSuccess
}
