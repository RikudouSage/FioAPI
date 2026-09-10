package fio

import "testing"

func TestOptionsAreAppliedInOrder(t *testing.T) {
	clientValue, err := NewClient("token",
		WithBaseURL("https://first.example/api"),
		WithBaseURL("https://second.example/rest/"),
		WithDebug(false),
		WithDebug(true),
	)
	if err != nil {
		t.Fatal(err)
	}
	got := clientValue.(*client)
	if got.baseURL != "https://second.example/rest" || !got.debug {
		t.Errorf("client options = baseURL %q, debug %v", got.baseURL, got.debug)
	}
}

func TestWithBaseURLRejectsMalformedURL(t *testing.T) {
	_, err := NewClient("token", WithBaseURL("https://example.com/%zz"))
	if err == nil {
		t.Fatal("NewClient() unexpectedly accepted a malformed URL")
	}
}

func TestNewClientUsesProductionURLByDefault(t *testing.T) {
	clientValue, err := NewClient("token")
	if err != nil {
		t.Fatal(err)
	}
	if got := clientValue.(*client).baseURL; got != "https://fioapi.fio.cz/v1/rest" {
		t.Errorf("baseURL = %q", got)
	}
}
