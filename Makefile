libfio.so libfio.h:
	go build -buildmode=c-shared -ldflags="-s -w" -o libfio.so ./cbindings

.PHONY: test test-c clean

test: test-c
	go test ./...

test-c: libfio.so libfio.h
	$(CC) -std=c11 -Wall -Wextra -Werror -Wno-unused-function -I. cbindings/tests/api_test.c -L. -lfio -o /tmp/fio-c-api-test
	LD_LIBRARY_PATH=. /tmp/fio-c-api-test

clean:
	rm -f libfio.so libfio.h
