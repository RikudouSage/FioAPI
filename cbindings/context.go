package main

/*
#include "headers/common.h"
*/
import "C"
import "context"

type contextHandle struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (receiver *contextHandle) Close() error {
	receiver.cancel()
	return nil
}

// FioNewContext creates a cancellable background context and writes its opaque
// handle to out. Closing the handle cancels the context.
//
//export FioNewContext
func FioNewContext(out *C.ContextHandle) C.FioResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioFailure
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctxHandle := &contextHandle{ctx: ctx, cancel: cancel}
	outHandle := registerHandle(ctxHandle)

	*out = C.ContextHandle(outHandle)

	clearLastError()
	return C.FioSuccess
}
