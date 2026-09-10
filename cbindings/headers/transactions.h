#pragma once

#include <stddef.h>
#include <stdint.h>

#include "common.h"

/**
 * A booked account transaction.
 *
 * date uses YYYY-MM-DD-HHMM format, including the numeric UTC offset supplied
 * by Fio. Nullable API values are represented by nullable pointers. Values
 * returned as part of a FioTransactions collection are owned by that
 * collection and must not be freed individually.
 */
typedef struct {
	int64_t id;
	char* date;
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
 * Account metadata and balances returned by FioGetAccountInfo.
 *
 * All strings are owned by this value and must be released together with it
 * by calling FioFreeAccountInfo. Nullable API values use nullable pointers.
 * date_start and date_end use YYYY-MM-DD-HHMM format.
 */
typedef struct {
	char* account_id;
	char* bank_id;
	char* currency;
	char* iban;
	char* bic;
	char* opening_balance;
	char* closing_balance;
	char* date_start;
	char* date_end;
	uint16_t* year_list;
	int64_t* id_list;
	int64_t* id_from;
	int64_t* id_to;
	int64_t* id_last_download;
} FioAccountInfo;

/**
 * A Czech domestic payment order.
 *
 * Required string pointers must be non-NULL. Optional string pointers may be
 * NULL and are borrowed only for the duration of FioIssueDomesticTransaction.
 * date uses YYYY-MM-DD format. A NULL or empty date, currency, or payment_type
 * receives the corresponding default defined by the Go API.
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
	const char* date;
	const char* message_for_recipient;
	const char* comment;
	const char* payment_type;
} FioDomesticTransaction;
