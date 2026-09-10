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

type TransactionValue[TType any] struct {
	Value TType
}

func (receiver TransactionValue[TType]) MarshalJSON() ([]byte, error) {
	val := value[TType]{
		Value: receiver.Value,
	}

	return json.Marshal(val)
}

func (receiver *TransactionValue[TType]) UnmarshalJSON(data []byte) error {
	var val value[TType]
	if err := json.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("failed unmarshalling inner struct: %w", err)
	}

	receiver.Value = val.Value
	return nil
}
