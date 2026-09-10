// Package fio provides a client for the Fio bank REST API.
//
// Create a [Client] with [NewClient], then use it to read account transactions,
// update the server-side last-pull marker, or issue domestic payments. All
// network operations accept a context so callers can control cancellation and
// timeouts.
//
// API tokens grant access to a bank account and should be treated as secrets.
// Refer to Fio bank's API documentation for the permissions required by each
// operation and for request frequency limits.
package fio
