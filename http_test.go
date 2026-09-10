package fio

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

type failingReadCloser struct{}

func (failingReadCloser) Read([]byte) (int, error) { return 0, errors.New("read failed") }
func (failingReadCloser) Close() error             { return nil }

func TestRequestMarshalsJSONBody(t *testing.T) {
	clientValue, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != `{"name":"value"}` {
			t.Errorf("body = %s", data)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	got, err := clientValue.(*client).request[struct {
		OK bool `json:"ok"`
	}](context.Background(), http.MethodPost, "/json", map[string]string{"name": "value"})
	if err != nil {
		t.Fatal(err)
	}
	if !got.OK {
		t.Error("response was not decoded")
	}
}

func TestRequestAcceptsReaderAndStringBodies(t *testing.T) {
	tests := []struct {
		name string
		body any
		want string
	}{
		{"reader", strings.NewReader("reader body"), "reader body"},
		{"string", "string body", "string body"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clientValue, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if string(data) != test.want {
					t.Errorf("body = %q, want %q", data, test.want)
				}
				w.WriteHeader(http.StatusNoContent)
			})
			if _, err := clientValue.(*client).request[any](context.Background(), http.MethodPost, "/body", test.body); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRequestReportsJSONMarshalError(t *testing.T) {
	instance := &client{baseURL: "https://example.invalid"}
	_, err := instance.request[any](context.Background(), http.MethodPost, "/", make(chan int))
	if err == nil || !strings.Contains(err.Error(), "failed marshalling body as JSON") {
		t.Fatalf("error = %v, want marshal error", err)
	}
}

func TestRequestRejectsWrongCustomDecoderType(t *testing.T) {
	clientValue, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("response"))
	})
	_, err := clientValue.(*client).request[string](context.Background(), http.MethodGet, "/custom", nil,
		withResponseBodyDecoder(func(io.Reader) (any, error) { return 42, nil }),
	)
	if err == nil || !strings.Contains(err.Error(), "expected TResult but got int") {
		t.Fatalf("error = %v, want decoder type error", err)
	}
}

func TestRequestReturnsCustomDecoderError(t *testing.T) {
	clientValue, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("response"))
	})
	want := errors.New("decoder failed")
	_, err := clientValue.(*client).request[string](context.Background(), http.MethodGet, "/custom", nil,
		withResponseBodyDecoder(func(io.Reader) (any, error) { return nil, want }),
	)
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want wrapped decoder error", err)
	}
}

func TestRequestReturnsConstructionError(t *testing.T) {
	instance := &client{baseURL: "https://example.invalid"}
	_, err := instance.request[any](context.Background(), "invalid\nmethod", "/", nil)
	if err == nil || !strings.Contains(err.Error(), "failed creating request") {
		t.Fatalf("error = %v, want request-construction error", err)
	}
}

func TestDebugErrorHandlesUnreadableResponseBody(t *testing.T) {
	previousClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       failingReadCloser{},
			Header:     make(http.Header),
		}, nil
	})}
	t.Cleanup(func() { http.DefaultClient = previousClient })

	instance := &client{baseURL: "https://example.invalid", debug: true}
	_, err := instance.request[any](context.Background(), http.MethodGet, "/", nil)
	if err == nil || !strings.Contains(err.Error(), "failed reading response body: read failed") {
		t.Fatalf("error = %v, want response-body read error", err)
	}
}

func TestHTTPOptionErrorIsReturned(t *testing.T) {
	instance := &client{baseURL: "https://example.invalid"}
	_, err := instance.request[any](context.Background(), http.MethodGet, "/", nil,
		func(*http.Request, *httpOptions) error { return io.ErrUnexpectedEOF },
	)
	if err == nil || !strings.Contains(err.Error(), io.ErrUnexpectedEOF.Error()) {
		t.Fatalf("error = %v, want option error", err)
	}
}
