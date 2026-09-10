package response

import "go.chrastecky.dev/fio-api/fio/dto"

type TransactionsResponse struct {
	AccountStatement struct {
		Info            dto.AccountInfo `json:"info"`
		TransactionList struct {
			Transactions []dto.Transaction `json:"transaction"`
		} `json:"transactionList"`
	} `json:"accountStatement"`
}
