package fio

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"go.chrastecky.dev/fio-api/fio/dto"
	"go.chrastecky.dev/fio-api/fio/internal/request"
	"go.chrastecky.dev/fio-api/fio/internal/response"
)

func (receiver *client) createXml(order request.PaymentImport) ([]byte, error) {
	out := bytes.NewBuffer(nil)

	order.ProvideDefaults()
	if _, err := io.WriteString(out, xml.Header); err != nil {
		return nil, fmt.Errorf("failed writing xml header: %w", err)
	}

	encoder := xml.NewEncoder(out)
	defer encoder.Close()
	if receiver.debug {
		encoder.Indent("", "  ")
	}
	if err := encoder.Encode(order); err != nil {
		return nil, fmt.Errorf("failed writing xml content: %w", err)
	}
	if err := encoder.Flush(); err != nil {
		return nil, fmt.Errorf("failed writing xml content: %w", err)
	}

	return out.Bytes(), nil
}

func (receiver *client) multipartImportBody(order request.PaymentImport) (content []byte, contentType string, err error) {
	out := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(out)
	defer writer.Close()

	if err := writer.WriteField("type", "xml"); err != nil {
		return nil, "", fmt.Errorf("failed writing multipart type field: %w", err)
	}
	if err := writer.WriteField("token", receiver.token); err != nil {
		return nil, "", fmt.Errorf("failed writing multipart token field: %w", err)
	}

	filePart, err := writer.CreateFormFile("file", "payments.xml")
	if err != nil {
		return nil, "", fmt.Errorf("failed creating multipart file part: %w", err)
	}

	xmlContent, err := receiver.createXml(order)
	if err != nil {
		return nil, "", fmt.Errorf("failed creating multipart file xml content: %w", err)
	}

	if _, err := filePart.Write(xmlContent); err != nil {
		return nil, "", fmt.Errorf("failed writing multipart file xml content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed closing multipart file writer: %w", err)
	}

	return out.Bytes(), writer.FormDataContentType(), nil
}

func (receiver *client) IssueDomesticTransactions(ctx context.Context, transaction dto.DomesticTransaction) error {
	order := request.PaymentImport{
		Orders: request.PaymentOrders{
			DomesticTransactions: []*dto.DomesticTransaction{&transaction},
		},
	}

	data, contentType, err := receiver.multipartImportBody(order)
	if err != nil {
		return fmt.Errorf("failed creating xml: %w", err)
	}
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

			return result, err
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
