# Fio API for Go

Go client and C bindings for reading transactions and issuing domestic payments
through the Fio banka REST API.

[API documentation](https://pkg.go.dev/go.chrastecky.dev/fio-api/fio)

## Installation

```sh
go get go.chrastecky.dev/fio-api/fio
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.chrastecky.dev/fio-api/fio"
)

func main() {
	client, err := fio.NewClient(os.Getenv("FIO_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	transactions, err := client.TransactionsSinceLastPull(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, transaction := range transactions {
		fmt.Printf("%d: %s %s\n",
			transaction.ID.Value,
			transaction.Amount.Value,
			transaction.Currency.Value,
		)
	}
}
```

## C bindings

Build `libfio.so` and its generated `libfio.h` header with:

```sh
make
```
