#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "libfio.h"

static int failures;

/* No error is represented by the size of an empty NUL-terminated string. */
static const size_t NO_ERROR = 1;

#define EXPECT(condition) do { \
	if (!(condition)) { \
		fprintf(stderr, "%s:%d: check failed: %s\n", __FILE__, __LINE__, #condition); \
		failures++; \
	} \
} while (0)

static void expect_error_contains(const char* expected) {
	size_t required = FioGetLastError(NULL, 0);
	EXPECT(required > 1);

	char* message = malloc(required);
	EXPECT(message != NULL);
	if (message == NULL) {
		return;
	}

	EXPECT(FioGetLastError(message, required) == required);
	EXPECT(strstr(message, expected) != NULL);

	if (required > 2) {
		char short_buffer[2] = {'x', 'x'};
		EXPECT(FioGetLastError(short_buffer, sizeof(short_buffer)) == required);
		EXPECT(short_buffer[1] == '\0');
	}
	free(message);
}

static void test_client_creation(void) {
	NewClientOptions options = {
		.token = "test-token",
		.clientUrl = NULL,
		.debug = false,
	};

	EXPECT(FioNewClient(NULL, options) == FioFailure);
	expect_error_contains("out is NULL");

	ClientHandle client = 0;
	options.token = NULL;
	EXPECT(FioNewClient(&client, options) == FioFailure);
	EXPECT(client == 0);
	expect_error_contains("token is NULL");

	options.token = "test-token";
	options.clientUrl = "not-a-url";
	EXPECT(FioNewClient(&client, options) == FioFailure);
	EXPECT(client == 0);
	expect_error_contains("missing scheme");

	options.clientUrl = "http://127.0.0.1:1";
	EXPECT(FioNewClient(&client, options) == FioSuccess);
	EXPECT(client != 0);
	EXPECT(FioGetLastError(NULL, 0) == NO_ERROR);
	EXPECT(FioCloseHandle(client) == FioSuccess);
	EXPECT(FioCloseHandle(client) == FioFailure);
	expect_error_contains("not registered");
}

static void test_context_and_handle_validation(void) {
	EXPECT(FioNewContext(NULL) == FioFailure);
	expect_error_contains("out is NULL");

	ContextHandle context = 0;
	EXPECT(FioNewContext(&context) == FioSuccess);
	EXPECT(context != 0);

	FioTransactions transactions = {
		.items = (FioTransaction*)1,
		.length = 99,
	};
	EXPECT(FioTransactionsSinceLastPull(999999, context, &transactions) == FioFailure);
	EXPECT(transactions.items == NULL);
	EXPECT(transactions.length == 0);
	expect_error_contains("handle 999999 not found");

	EXPECT(FioTransactionsSinceLastPull(0, context, NULL) == FioFailure);
	expect_error_contains("out is NULL");

	EXPECT(FioCloseHandle(context) == FioSuccess);
	EXPECT(FioCloseHandle(context) == FioFailure);
}

static void test_argument_conversion(void) {
	NewClientOptions options = {
		.token = "test-token",
		.clientUrl = "http://127.0.0.1:1",
		.debug = false,
	};
	ClientHandle client = 0;
	ContextHandle context = 0;
	EXPECT(FioNewClient(&client, options) == FioSuccess);
	EXPECT(FioNewContext(&context) == FioSuccess);

	FioTransactions transactions = {0};
	EXPECT(FioTransactionsByDate(client, context, NULL, "2025-01-02", &transactions) == FioFailure);
	expect_error_contains("start_date is NULL");
	EXPECT(FioTransactionsByDate(client, context, "invalid", "2025-01-02", &transactions) == FioFailure);
	expect_error_contains("invalid start_date");
	EXPECT(FioSetLastFailedTransactionDate(client, context, NULL) == FioFailure);
	expect_error_contains("date is NULL");
	EXPECT(FioSetLastFailedTransactionDate(client, context, "2025-02-30") == FioFailure);
	expect_error_contains("invalid date");

	FioDomesticTransaction payment = {0};
	EXPECT(FioIssueDomesticTransaction(client, context, payment) == FioFailure);
	expect_error_contains("transaction.account_from is NULL");

	payment.account_from = "123";
	EXPECT(FioIssueDomesticTransaction(client, context, payment) == FioFailure);
	expect_error_contains("transaction.amount is NULL");

	payment.amount = "not-an-amount";
	EXPECT(FioIssueDomesticTransaction(client, context, payment) == FioFailure);
	expect_error_contains("invalid transaction.amount");

	FioFreeTransactions(NULL);
	FioFreeTransactions(&transactions);
	EXPECT(transactions.items == NULL);
	EXPECT(transactions.length == 0);

	EXPECT(FioCloseHandle(context) == FioSuccess);
	EXPECT(FioCloseHandle(client) == FioSuccess);
}

int main(void) {
	test_client_creation();
	test_context_and_handle_validation();
	test_argument_conversion();

	if (failures != 0) {
		fprintf(stderr, "%d C API test(s) failed\n", failures);
		return EXIT_FAILURE;
	}

	puts("C API tests passed");
	return EXIT_SUCCESS;
}
