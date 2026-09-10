package main

/*
#include "headers/transactions.h"
*/
import "C"

//export FioIssueDomesticTransaction
func FioIssueDomesticTransaction(client C.ClientHandle, ctx C.ContextHandle, transaction C.FioDomesticTransaction) C.FioResult {
	clientGo, ctxGo, err := getCommonHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	transactionGo, err := domesticTransactionFromC(transaction)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	if err = clientGo.IssueDomesticTransaction(ctxGo, transactionGo); err != nil {
		setLastError(err)
		return C.FioFailure
	}

	clearLastError()
	return C.FioSuccess
}
