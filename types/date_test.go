package types

import (
	"testing"
	"time"
)

func TestDateTextRoundTrip(t *testing.T) {
	want := "2025-06-07"
	var date Date
	if err := date.UnmarshalText([]byte(want)); err != nil {
		t.Fatal(err)
	}
	text, err := date.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if string(text) != want {
		t.Errorf("round trip = %q, want %q", text, want)
	}
	if _, offset := date.AsTime().Zone(); offset != 0 {
		t.Errorf("parsed date offset = %d, want UTC", offset)
	}
}

func TestDateRejectsInvalidText(t *testing.T) {
	var date Date
	if err := date.UnmarshalText([]byte("2025-02-30")); err == nil {
		t.Fatal("UnmarshalText() unexpectedly accepted an invalid date")
	}
}

func TestDateValueEqualsString(t *testing.T) {
	date := Date(time.Date(2025, 6, 7, 0, 0, 0, 0, time.UTC))
	got, err := date.Value()
	if err != nil {
		t.Fatal(err)
	}
	if got != date.String() {
		t.Errorf("Value() = %q, want String() = %q", got, date.String())
	}
}

func TestDateScanEqualsString(t *testing.T) {
	want := Date(time.Date(2025, 6, 7, 0, 0, 0, 0, time.UTC))
	var got Date
	if err := got.Scan(want.String()); err != nil {
		t.Fatal(err)
	}
	if got.String() != want.String() {
		t.Errorf("Scan() produced %q, want %q", got.String(), want.String())
	}
}

func TestTimezonedDateTextRoundTrip(t *testing.T) {
	want := "2025-06-07+0230"
	var date TimezonedDate
	if err := date.UnmarshalText([]byte(want)); err != nil {
		t.Fatal(err)
	}
	if got := date.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	wantOffset := int((2*time.Hour + 30*time.Minute) / time.Second)
	if _, offset := date.AsTime().Zone(); offset != wantOffset {
		t.Errorf("parsed offset = %s, want 2h30m", time.Duration(offset)*time.Second)
	}
	text, err := date.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if string(text) != want {
		t.Errorf("MarshalText() = %q, want %q", text, want)
	}
}

func TestTimezonedDateRejectsInvalidText(t *testing.T) {
	var date TimezonedDate
	if err := date.UnmarshalText([]byte("2025-06-07")); err == nil {
		t.Fatal("UnmarshalText() unexpectedly accepted a date without an offset")
	}
}

func TestTimezonedDateValueEqualsString(t *testing.T) {
	date := TimezonedDate(time.Date(2025, 6, 7, 0, 0, 0, 0, time.FixedZone("test", 2*60*60)))
	got, err := date.Value()
	if err != nil {
		t.Fatal(err)
	}
	if got != date.String() {
		t.Errorf("Value() = %q, want String() = %q", got, date.String())
	}
}

func TestTimezonedDateScanEqualsString(t *testing.T) {
	want := TimezonedDate(time.Date(2025, 6, 7, 0, 0, 0, 0, time.FixedZone("test", 2*60*60)))
	var got TimezonedDate
	if err := got.Scan(want.String()); err != nil {
		t.Fatal(err)
	}
	if got.String() != want.String() {
		t.Errorf("Scan() produced %q, want %q", got.String(), want.String())
	}
}
