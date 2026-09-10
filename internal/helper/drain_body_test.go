package helper

import (
	"io"
	"strings"
	"testing"
)

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error {
	b.closed = true
	return nil
}

func TestDrainBodyReadsAndCloses(t *testing.T) {
	body := &trackedBody{Reader: strings.NewReader("remaining response")}
	DrainBody(body)
	if !body.closed {
		t.Error("body was not closed")
	}
	remaining, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(remaining) != 0 {
		t.Errorf("body still contains %q", remaining)
	}
}
