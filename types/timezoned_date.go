package types

import (
	"fmt"
	"time"
)

// TimezonedDate represents a calendar date with a numeric UTC offset. Its text
// format is YYYY-MM-DD-HHMM, where the final component is the UTC offset.
type TimezonedDate time.Time

// AsTime converts the value to time.Time.
func (receiver TimezonedDate) AsTime() time.Time {
	return time.Time(receiver)
}

// MarshalText implements encoding.TextMarshaler using the Fio API date format.
func (receiver TimezonedDate) MarshalText() (text []byte, err error) {
	return []byte(receiver.AsTime().Format("2006-01-02-0700")), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for Fio API date values.
func (receiver *TimezonedDate) UnmarshalText(text []byte) error {
	parsed, err := time.Parse("2006-01-02-0700", string(text))
	if err != nil {
		return fmt.Errorf("failed parsing date: %w", err)
	}

	*receiver = TimezonedDate(parsed)
	return nil
}
