package types

import (
	"fmt"
	"time"
)

// Date represents a calendar date without a time of day. Its text format is
// YYYY-MM-DD.
type Date time.Time

// AsTime converts the date to time.Time.
func (receiver Date) AsTime() time.Time {
	return time.Time(receiver)
}

// MarshalText implements encoding.TextMarshaler using YYYY-MM-DD format.
func (receiver Date) MarshalText() (text []byte, err error) {
	return []byte(receiver.AsTime().Format("2006-01-02")), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for YYYY-MM-DD values.
func (receiver *Date) UnmarshalText(text []byte) error {
	parsed, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		return fmt.Errorf("failed parsing date: %w", err)
	}

	*receiver = Date(parsed)
	return nil
}
