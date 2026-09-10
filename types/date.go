package types

import (
	"fmt"
	"time"
)

type Date time.Time

func (receiver Date) AsTime() time.Time {
	return time.Time(receiver)
}

func (receiver Date) MarshalText() (text []byte, err error) {
	return []byte(receiver.AsTime().Format("2006-01-02")), nil
}

func (receiver *Date) UnmarshalText(text []byte) error {
	parsed, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		return fmt.Errorf("failed parsing date: %w", err)
	}

	*receiver = Date(parsed)
	return nil
}
