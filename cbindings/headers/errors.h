#pragma once

#include <stdlib.h>
#include <string.h>

/* Thread-local error storage used by the exported C API. */
static __thread char* fio_last_error;

/* Clears the calling thread's stored error. */
static void fio_clear_last_error(void) {
	if (fio_last_error) {
		free(fio_last_error);
		fio_last_error = NULL;
	}
}

/* Replaces the calling thread's stored error with a copy of msg. */
static void fio_set_last_error_copy(const char* msg) {
	fio_clear_last_error();
	if (!msg) {
		return;
	}

	size_t len = strlen(msg);
	char* copy = (char*)malloc(len + 1);
	if (!copy) {
		return;
	}
	memcpy(copy, msg, len + 1);
	fio_last_error = copy;
}

/*
 * Copies the calling thread's error to buf and returns the required size,
 * including the terminating NUL byte. A return value of 1 denotes no error.
 */
static size_t fio_get_last_error(char* buf, size_t buf_len) {
	if (!fio_last_error) {
		if (buf && buf_len > 0) {
			buf[0] = '\0';
		}
		return 1;
	}

	size_t len = strlen(fio_last_error) + 1;
	if (buf && buf_len > 0) {
		size_t to_copy = len <= buf_len ? len : buf_len - 1;
		memcpy(buf, fio_last_error, to_copy);
		buf[to_copy] = '\0';
	}
	return len;
}
