#pragma once

#include <stddef.h>
#include <stdint.h>

typedef enum {
	FioSuccess = 0,
	FioFailure = 1,
} FioResult;

typedef uint64_t Handle;
typedef Handle ClientHandle;
typedef Handle ContextHandle;