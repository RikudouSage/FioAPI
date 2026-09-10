#pragma once

#include <stddef.h>
#include <stdint.h>

/** Result returned by C API operations. */
typedef enum {
	/** The operation completed successfully. */
	FioSuccess = 0,
	/** The operation failed; use FioGetLastError to obtain details. */
	FioFailure = 1,
} FioResult;

/** Opaque handle to a Go-managed object. */
typedef uint64_t Handle;
/** Opaque handle to a Fio API client. */
typedef Handle ClientHandle;
/** Opaque handle to a cancellable operation context. */
typedef Handle ContextHandle;
