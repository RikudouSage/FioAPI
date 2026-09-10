package main

/*
#include "headers/common.h"
*/
import "C"

//export FioCloseHandle
func FioCloseHandle(handleID C.Handle) C.FioResult {
	err := unregisterHandle(handle(handleID))
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	clearLastError()
	return C.FioSuccess
}
