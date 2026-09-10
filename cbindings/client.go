package main

/*
#include <stdbool.h>
#include "headers/common.h"

#ifndef FIO_NEW_CLIENT_OPTIONS
#define FIO_NEW_CLIENT_OPTIONS
// Options used to create a Fio API client.
typedef struct {
	// Required Fio API token. The string is borrowed for the call.
	const char* token;
	// Optional API base URL. NULL selects the production endpoint.
	const char* clientUrl;
	// Enables additional response details in errors and formatted request XML.
	bool debug;
} NewClientOptions;
#endif

*/
import "C"
import "go.chrastecky.dev/fio-api/fio"

// FioNewClient creates a client and writes its opaque handle to out.
// The caller must release the handle with FioCloseHandle.
//
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
