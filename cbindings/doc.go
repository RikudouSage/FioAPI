// Package main implements C bindings for the Fio bank API client.
//
// The package is built as a C shared library. Its exported functions use opaque
// handles for Go-managed clients and contexts; callers must release those
// handles with FioCloseHandle. Functions returning FioResult store diagnostic
// details in thread-local error state when they return FioFailure.
package main
