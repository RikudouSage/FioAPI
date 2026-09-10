#pragma once

#include <stddef.h>
#include <stdint.h>

#include "common.h"

/**
 * A booked account transaction.
 *
 * date_ms is expressed as milliseconds since the Unix epoch. Nullable API
 * values are represented by nullable pointers. Values returned as part of a
 * FioTransactions collection are owned by that collection and must not be
 * freed individually.
 */
typedef struct {
	int64_t id;
	uint64_t date_ms;
	char* amount;
	char* currency;
	char* counterparty_account;
	char* counterparty_name;
	char* counterparty_bank_code;
	char* counterparty_bank_name;
	char* constant_symbol;
	char* variable_symbol;
	char* specific_symbol;
	char* user_identity;
	char* transaction_type;
	char* performed_by;
	char* additional_info;
	char* comment;
	char* bic;
	int64_t* instruction_id;
	char* payer_reference;
} FioTransaction;

/**
 * An owned collection of booked transactions.
 * Release it with FioFreeTransactions, including when length is zero.
 */
typedef struct {
	FioTransaction* items;
	size_t length;
} FioTransactions;

/**
 * A Czech domestic payment order.
 *
 * Required string pointers must be non-NULL. Optional string pointers may be
 * NULL and are borrowed only for the duration of FioIssueDomesticTransaction.
 * date_ms is expressed as milliseconds since the Unix epoch. Empty currency
 * and payment_type strings receive the defaults defined by the Go API.
 */
typedef struct {
	const char* account_from;
	const char* currency;
	const char* amount;
	const char* account_to;
	const char* bank_code;
	const char* constant_symbol;
	const char* variable_symbol;
	const char* specific_symbol;
	uint64_t date_ms;
	const char* message_for_recipient;
	const char* comment;
	const char* payment_type;
} FioDomesticTransaction;
