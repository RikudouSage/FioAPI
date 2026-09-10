package types

import "testing"

func transactionType(value string) TransactionType {
	return TransactionType{innerType: value}
}

func TestTransactionTypeClassifiers(t *testing.T) {
	tests := []struct {
		name      string
		typ       TransactionType
		predicate func(TransactionType) bool
	}{
		{"incoming", transactionType(typeIncoming), TransactionType.IsIncoming},
		{"outgoing", transactionType(typePayment), TransactionType.IsOutgoing},
		{"immediate", transactionType(typeInstantOutgoing), TransactionType.IsImmediate},
		{"internal", transactionType(typeTransferInsideAccount), TransactionType.IsInternal},
		{"card", transactionType(typeCardPayment), TransactionType.IsCardPayment},
		{"cash", transactionType(typeCashDeposit), TransactionType.IsCash},
		{"direct debit", transactionType(typeDirectDebit), TransactionType.IsDirectDebit},
		{"fee", transactionType(typeRecordedFee), TransactionType.IsFee},
		{"interest", transactionType(typeInterestCredited), TransactionType.IsInterest},
		{"correction", transactionType(typeCorrection), TransactionType.IsCorrection},
		{"known", transactionType(typeInstantEuroIncoming), TransactionType.IsKnown},
	}
	unknown := transactionType("unknown")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.predicate(test.typ) {
				t.Errorf("predicate returned false for %q", test.typ)
			}
			if test.predicate(unknown) {
				t.Error("predicate returned true for an unknown transaction type")
			}
		})
	}
}

func TestTransactionTypeMarshalText(t *testing.T) {
	typ := transactionType(typeCardPayment)
	text, err := typ.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if string(text) != typeCardPayment {
		t.Errorf("MarshalText() = %q", text)
	}
}

func TestTransactionTypeValueEqualsString(t *testing.T) {
	typ := transactionType(typeCardPayment)
	got, err := typ.Value()
	if err != nil {
		t.Fatal(err)
	}
	if got != typ.String() {
		t.Errorf("Value() = %q, want String() = %q", got, typ.String())
	}
}
