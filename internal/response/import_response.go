package response

import "encoding/xml"

type ImportResultErrorCode int

const (
	ImportResultSuccess         ImportResultErrorCode = 0
	ImportResultOrderError      ImportResultErrorCode = 1
	ImportResultOrderWarning    ImportResultErrorCode = 2
	ImportResultSyntacticError  ImportResultErrorCode = 11
	ImportResultRequestEmpty    ImportResultErrorCode = 12
	ImportResultRequestTooLarge ImportResultErrorCode = 13
	ImportResultRequestEmpty2   ImportResultErrorCode = 14
)

type ImportResultStatus string

const (
	ImportResultOK      ImportResultStatus = "ok"
	ImportResultError   ImportResultStatus = "error"
	ImportResultWarning ImportResultStatus = "warning"
)

type ImportResponse struct {
	XMLName xml.Name `xml:"responseImport"`

	Result        ImportResult  `xml:"result"`
	OrdersDetails OrdersDetails `xml:"ordersDetails"`
}

type ImportResult struct {
	ErrorCode     ImportResultErrorCode `xml:"errorCode"`
	IDInstruction int64                 `xml:"idInstruction"`
	Status        ImportResultStatus    `xml:"status"`
	Sums          ImportSums            `xml:"sums"`
}

type ImportSums struct {
	Sums []ImportSum `xml:"sum"`
}

type ImportSum struct {
	Currency  string `xml:"id,attr"`
	SumCredit string `xml:"sumCredit"`
	SumDebit  string `xml:"sumDebet"`
}

type OrdersDetails struct {
	Details []OrderDetail `xml:"detail"`
}

type OrderDetail struct {
	ID       int            `xml:"id,attr"`
	Messages []OrderMessage `xml:"messages>message"`
}

type OrderMessage struct {
	ErrorCode int    `xml:"errorCode,attr"`
	Status    string `xml:"status,attr"`
	Text      string `xml:",chardata"`
}
