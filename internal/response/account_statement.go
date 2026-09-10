package response

import (
	"github.com/shopspring/decimal"
	"go.chrastecky.dev/fio-api/fio/types"
)

type AccountInfo struct {
	AccountID      string              `json:"accountId"`
	BankID         string              `json:"bankId"`
	Currency       string              `json:"currency"`
	IBAN           string              `json:"iban"`
	BIC            string              `json:"bic"`
	OpeningBalance decimal.Decimal     `json:"openingBalance"`
	ClosingBalance decimal.Decimal     `json:"closingBalance"`
	DateStart      types.TimezonedDate `json:"dateStart"`
	DateEnd        types.TimezonedDate `json:"dateEnd"`
	YearList       *uint16             `json:"yearList"`
	IDList         *int64              `json:"idList"`
	IDFrom         *int64              `json:"idFrom"`
	IDTo           *int64              `json:"idTo"`
	IDLastDownload *int64              `json:"idLastDownload"`
}
