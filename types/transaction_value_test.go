package types

import (
	"encoding/json"
	"testing"
)

func TestTransactionValueJSON(t *testing.T) {
	var got TransactionValue[*string]
	if err := json.Unmarshal([]byte(`{"value":"reference","name":"VS","id":5}`), &got); err != nil {
		t.Fatal(err)
	}
	if got.Value == nil || *got.Value != "reference" {
		t.Fatalf("decoded value = %v, want reference", got.Value)
	}

	data, err := json.Marshal(TransactionValue[int]{Value: 7})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"value":7,"name":"","id":0}` {
		t.Errorf("MarshalJSON() = %s", data)
	}
}

func TestUnknownTransactionTypeIsPreserved(t *testing.T) {
	var typ TransactionType
	if err := typ.UnmarshalText([]byte("new bank category")); err != nil {
		t.Fatal(err)
	}
	if typ.String() != "new bank category" || typ.IsKnown() {
		t.Errorf("unknown type = %q, IsKnown = %v", typ, typ.IsKnown())
	}
}
