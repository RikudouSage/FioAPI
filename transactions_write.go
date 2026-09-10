package fio

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/samber/lo"
	"go.chrastecky.dev/fio-api/fio/dto"
	"go.chrastecky.dev/fio-api/fio/internal/request"
	"go.chrastecky.dev/fio-api/fio/internal/response"
)

func (receiver *client) createXml(order request.PaymentImport) []byte {
	out := bytes.NewBuffer(nil)

	order.ProvideDefaults()
	lo.Must(io.WriteString(out, xml.Header))

	encoder := xml.NewEncoder(out)
	if receiver.debug {
		encoder.Indent("", "  ")
	}
	lo.Must0(encoder.Encode(order))
	lo.Must0(encoder.Close())

	return out.Bytes()
}

func (receiver *client) multipartImportBody(order request.PaymentImport) (content []byte, contentType string) {
	out := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(out)

	lo.Must0(writer.WriteField("type", "xml"))
	lo.Must0(writer.WriteField("token", receiver.token))

	filePart := lo.Must(writer.CreateFormFile("file", "payments.xml"))

	lo.Must(filePart.Write(receiver.createXml(order)))
	lo.Must0(writer.Close())

	return out.Bytes(), writer.FormDataContentType()
}

func (receiver *client) IssueDomesticTransaction(ctx context.Context, transaction dto.DomesticTransaction) error {
	order := request.PaymentImport{
		Orders: request.PaymentOrders{
			DomesticTransactions: []*dto.DomesticTransaction{&transaction},
		},
	}

	data, contentType := receiver.multipartImportBody(order)
	resp, err := receiver.request[response.ImportResponse](
		ctx,
		http.MethodPost,
		"/import/",
		data,
		withHeader("Content-Type", contentType),
		withResponseBodyDecoder(func(body io.Reader) (any, error) {
			var result response.ImportResponse
			decoder := xml.NewDecoder(body)

			if err := decoder.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed decoding xml: %w", err)
			}

			return result, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("failed sending request: %w", err)
	}
	if resp.Result.Status != response.ImportResultOK {
		return fmt.Errorf("failed sending request, status %s (code: %d)", resp.Result.Status, resp.Result.ErrorCode)
	}

	return nil
}
