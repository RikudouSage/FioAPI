package types

import (
	"fmt"
	"time"
)

type FioDate time.Time

func (receiver FioDate) AsTime() time.Time {
	return time.Time(receiver)
}

func (receiver FioDate) MarshalText() (text []byte, err error) {
	return []byte(receiver.AsTime().Format("2006-01-02-0700")), nil
}

func (receiver *FioDate) UnmarshalText(text []byte) error {
	parsed, err := time.Parse("2006-01-02-0700", string(text))
	if err != nil {
		return fmt.Errorf("failed parsing date: %w", err)
	}

	*receiver = FioDate(parsed)
	return nil
}
