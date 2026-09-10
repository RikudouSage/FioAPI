package main

/*
#include <stdbool.h>
#include "headers/common.h"

#ifndef FIO_NEW_CLIENT_OPTIONS
#define FIO_NEW_CLIENT_OPTIONS
typedef struct {
	const char* token;
	const char* clientUrl;
	bool debug;
} NewClientOptions;
#endif

*/
import "C"
import "go.chrastecky.dev/fio-api/fio"

//export FioNewClient
func FioNewClient(out *C.ClientHandle, options C.NewClientOptions) C.FioResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioFailure
	}

	if options.token == nil {
		setLastError(nullPointerError("token"))
		return C.FioFailure
	}

	goOptions := make([]fio.Option, 0, 2)
	goOptions = append(goOptions, fio.WithDebug(bool(options.debug)))
	if options.clientUrl != nil {
		goOptions = append(goOptions, fio.WithBaseURL(C.GoString(options.clientUrl)))
	}

	client, err := fio.NewClient(C.GoString(options.token), goOptions...)
	if err != nil {
		setLastError(err)
		return C.FioFailure
	}

	result := registerHandle(client)
	*out = C.ClientHandle(result)

	clearLastError()
	return C.FioSuccess
}
