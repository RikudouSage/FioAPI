package main

import (
	"errors"
	"testing"
)

type testCloser struct {
	closed bool
	err    error
}

func (c *testCloser) Close() error {
	c.closed = true
	return c.err
}

func TestHandleLifecycle(t *testing.T) {
	value := &testCloser{}
	id := registerHandle(value)
	t.Cleanup(func() {
		// Leave the global registry clean if the test fails before unregistering.
		handlesMutex.Lock()
		delete(handles, id)
		handlesMutex.Unlock()
	})

	got, err := getHandleObj[*testCloser](id)
	if err != nil {
		t.Fatalf("getHandleObj() error = %v", err)
	}
	if got != value {
		t.Fatalf("getHandleObj() = %p, want %p", got, value)
	}
	if _, err := getHandleObj[string](id); err == nil {
		t.Fatal("getHandleObj() unexpectedly accepted the wrong type")
	}
	if err := unregisterHandle(id); err != nil {
		t.Fatalf("unregisterHandle() error = %v", err)
	}
	if !value.closed {
		t.Error("unregisterHandle() did not close the registered object")
	}
	if _, err := getHandleObj[*testCloser](id); err == nil {
		t.Fatal("getHandleObj() found an unregistered handle")
	}
	if err := unregisterHandle(id); err == nil {
		t.Fatal("unregisterHandle() unexpectedly accepted an unknown handle")
	}
}

func TestUnregisterHandleIgnoresCloseError(t *testing.T) {
	value := &testCloser{err: errors.New("close failed")}
	id := registerHandle(value)
	if err := unregisterHandle(id); err != nil {
		t.Fatalf("unregisterHandle() error = %v", err)
	}
	if !value.closed {
		t.Error("registered object was not closed")
	}
}
