package types

import (
	"encoding/json"
	"fmt"
)

type value[TType any] struct {
	Value TType  `json:"value"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
}

// TransactionValue contains a transaction field's value. It adapts Fio's
// metadata-bearing column object to a compact Go value.
type TransactionValue[TType any] struct {
	// Value is the decoded transaction field value.
	Value TType
}

// MarshalJSON implements json.Marshaler using Fio's column object format.
func (receiver TransactionValue[TType]) MarshalJSON() ([]byte, error) {
	val := value[TType]{
		Value: receiver.Value,
	}

	return json.Marshal(val)
}

// UnmarshalJSON implements json.Unmarshaler for Fio's column object format.
// Column metadata is intentionally discarded.
func (receiver *TransactionValue[TType]) UnmarshalJSON(data []byte) error {
	var val value[TType]
	if err := json.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("failed unmarshalling inner struct: %w", err)
	}

	receiver.Value = val.Value
	return nil
}
