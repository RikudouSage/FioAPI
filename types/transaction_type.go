package types

import (
	"database/sql/driver"
	"fmt"
)

const (
	TypeIncomingInsideBank               = "Příjem převodem uvnitř banky"
	TypeOutgoingInsideBank               = "Platba převodem uvnitř banky"
	TypeCashDepositAtCounter             = "Vklad pokladnou"
	TypeCashWithdrawalAtCounter          = "Výběr pokladnou"
	TypeCashDeposit                      = "Vklad v hotovosti"
	TypeCashWithdrawal                   = "Výběr v hotovosti"
	TypePayment                          = "Platba"
	TypeIncoming                         = "Příjem"
	TypeCashlessPayment                  = "Bezhotovostní platba"
	TypeCashlessIncoming                 = "Bezhotovostní příjem"
	TypeCardPayment                      = "Platba kartou"
	TypeLoanInterest                     = "Úrok z úvěru"
	TypePenaltyFee                       = "Sankční poplatek"
	TypeCourierOutgoing                  = "Posel – předání"
	TypeCourierIncoming                  = "Posel – příjem"
	TypeTransferInsideAccount            = "Převod uvnitř konta"
	TypeInterestCredited                 = "Připsaný úrok"
	TypeInterestPaid                     = "Vyplacený úrok"
	TypeInterestTax                      = "Odvod daně z úroků"
	TypeRecordedInterest                 = "Evidovaný úrok"
	TypeFee                              = "Poplatek"
	TypeRecordedFee                      = "Evidovaný poplatek"
	TypeInterBankAccountTransferOutgoing = "Převod mezi bankovními konty (platba)"
	TypeInterBankAccountTransferIncoming = "Převod mezi bankovními konty (příjem)"
	TypeUnidentifiedBankAccountPayment   = "Neidentifikovaná platba z bankovního konta"
	TypeUnidentifiedBankAccountIncoming  = "Neidentifikovaný příjem na bankovní konto"
	TypeOwnBankAccountPayment            = "Vlastní platba z bankovního konta"
	TypeOwnBankAccountIncoming           = "Vlastní příjem na bankovní konto"
	TypeOwnCounterPayment                = "Vlastní platba pokladnou"
	TypeOwnCounterIncoming               = "Vlastní příjem pokladnou"
	TypeCorrection                       = "Opravný pohyb"
	TypeFeeReceived                      = "Přijatý poplatek"
	TypeForeignCurrencyPayment           = "Platba v jiné měně"
	TypeCardFee                          = "Poplatek – platební karta"
	TypeDirectDebit                      = "Inkaso"
	TypeDirectDebitIncoming              = "Inkaso ve prospěch účtu"
	TypeDirectDebitOutgoing              = "Inkaso z účtu"
	TypeDirectDebitIncomingOtherBank     = "Příjem inkasa z cizí banky"
	TypeInstantIncoming                  = "Okamžitá příchozí platba"
	TypeInstantOutgoing                  = "Okamžitá odchozí platba"
	TypeMortgageInsuranceFee             = "Poplatek - pojištění hypotéky"
	TypeInstantEuroIncoming              = "Okamžitá příchozí Europlatba"
	TypeInstantEuroOutgoing              = "Okamžitá odchozí Europlatba"
)

// TransactionType is the localized transaction category reported by Fio
// bank. Unknown values are preserved and can be detected with [TransactionType.IsKnown].
type TransactionType struct {
	innerType string
}

func NewTransactionType(raw string) TransactionType {
	return TransactionType{
		innerType: raw,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (receiver TransactionType) MarshalText() ([]byte, error) {
	return []byte(receiver.innerType), nil
}

// UnmarshalText implements encoding.TextUnmarshaler. It accepts unknown values
// so newer server-side categories remain available to callers.
func (receiver *TransactionType) UnmarshalText(text []byte) error {
	*receiver = TransactionType{innerType: string(text)}
	return nil
}

// String returns the category text supplied by Fio bank.
func (receiver TransactionType) String() string {
	return receiver.innerType
}

// Value implements [driver.Valuer] by returning the transaction Type's
// original string representation.
func (receiver TransactionType) Value() (driver.Value, error) {
	return receiver.String(), nil
}

// Scan implements database/sql.Scanner for transaction Type strings.
func (receiver *TransactionType) Scan(src any) error {
	if str, ok := src.(string); ok {
		return receiver.UnmarshalText([]byte(str))
	}

	return fmt.Errorf("the value must be a string, %T given", src)
}

// IsIncoming reports whether the Type represents incoming funds.
func (receiver TransactionType) IsIncoming() bool {
	switch receiver.innerType {
	case
		TypeIncomingInsideBank,
		TypeIncoming,
		TypeCashlessIncoming,
		TypeCourierIncoming,
		TypeInterBankAccountTransferIncoming,
		TypeUnidentifiedBankAccountIncoming,
		TypeOwnBankAccountIncoming,
		TypeOwnCounterIncoming,
		TypeFeeReceived,
		TypeDirectDebitIncoming,
		TypeDirectDebitIncomingOtherBank,
		TypeInstantIncoming,
		TypeInstantEuroIncoming:
		return true
	default:
		return false
	}
}

// IsOutgoing reports whether the Type represents outgoing funds.
func (receiver TransactionType) IsOutgoing() bool {
	switch receiver.innerType {
	case
		TypeOutgoingInsideBank,
		TypeCashWithdrawalAtCounter,
		TypeCashWithdrawal,
		TypePayment,
		TypeCashlessPayment,
		TypeCardPayment,
		TypeLoanInterest,
		TypePenaltyFee,
		TypeCourierOutgoing,
		TypeInterBankAccountTransferOutgoing,
		TypeUnidentifiedBankAccountPayment,
		TypeOwnBankAccountPayment,
		TypeOwnCounterPayment,
		TypeFee,
		TypeCardFee,
		TypeDirectDebitOutgoing,
		TypeMortgageInsuranceFee,
		TypeForeignCurrencyPayment,
		TypeInstantOutgoing,
		TypeInstantEuroOutgoing:
		return true
	default:
		return false
	}
}

// IsImmediate reports whether the Type represents an instant payment.
func (receiver TransactionType) IsImmediate() bool {
	switch receiver.innerType {
	case
		TypeInstantIncoming,
		TypeInstantOutgoing,
		TypeInstantEuroIncoming,
		TypeInstantEuroOutgoing:
		return true
	default:
		return false
	}
}

// IsInternal reports whether the Type represents a transfer within Fio bank
// or between the account owner's accounts.
func (receiver TransactionType) IsInternal() bool {
	switch receiver.innerType {
	case
		TypeIncomingInsideBank,
		TypeOutgoingInsideBank,
		TypeTransferInsideAccount,
		TypeOwnBankAccountPayment,
		TypeOwnBankAccountIncoming:
		return true
	default:
		return false
	}
}

// IsCardPayment reports whether the Type represents a card payment.
func (receiver TransactionType) IsCardPayment() bool {
	return receiver.innerType == TypeCardPayment
}

// IsCash reports whether the Type represents a cash transaction.
func (receiver TransactionType) IsCash() bool {
	switch receiver.innerType {
	case
		TypeCashDepositAtCounter,
		TypeCashWithdrawalAtCounter,
		TypeCashDeposit,
		TypeCashWithdrawal,
		TypeOwnCounterPayment,
		TypeOwnCounterIncoming:
		return true
	default:
		return false
	}
}

// IsDirectDebit reports whether the Type represents a direct debit.
func (receiver TransactionType) IsDirectDebit() bool {
	switch receiver.innerType {
	case
		TypeDirectDebit,
		TypeDirectDebitIncoming,
		TypeDirectDebitOutgoing,
		TypeDirectDebitIncomingOtherBank:
		return true
	default:
		return false
	}
}

// IsFee reports whether the Type represents a fee.
func (receiver TransactionType) IsFee() bool {
	switch receiver.innerType {
	case
		TypePenaltyFee,
		TypeFee,
		TypeRecordedFee,
		TypeFeeReceived,
		TypeCardFee,
		TypeMortgageInsuranceFee:
		return true
	default:
		return false
	}
}

// IsInterest reports whether the Type represents interest or interest tax.
func (receiver TransactionType) IsInterest() bool {
	switch receiver.innerType {
	case
		TypeLoanInterest,
		TypeInterestCredited,
		TypeInterestPaid,
		TypeInterestTax,
		TypeRecordedInterest:
		return true
	default:
		return false
	}
}

// IsCorrection reports whether the Type represents a corrective entry.
func (receiver TransactionType) IsCorrection() bool {
	return receiver.innerType == TypeCorrection
}

// IsKnown reports whether the Type is recognized by this version of the
// package.
func (receiver TransactionType) IsKnown() bool {
	switch receiver.innerType {
	case
		TypeIncomingInsideBank,
		TypeOutgoingInsideBank,
		TypeCashDepositAtCounter,
		TypeCashWithdrawalAtCounter,
		TypeCashDeposit,
		TypeCashWithdrawal,
		TypePayment,
		TypeIncoming,
		TypeCashlessPayment,
		TypeCashlessIncoming,
		TypeCardPayment,
		TypeLoanInterest,
		TypePenaltyFee,
		TypeCourierOutgoing,
		TypeCourierIncoming,
		TypeTransferInsideAccount,
		TypeInterestCredited,
		TypeInterestPaid,
		TypeInterestTax,
		TypeRecordedInterest,
		TypeFee,
		TypeRecordedFee,
		TypeInterBankAccountTransferOutgoing,
		TypeInterBankAccountTransferIncoming,
		TypeUnidentifiedBankAccountPayment,
		TypeUnidentifiedBankAccountIncoming,
		TypeOwnBankAccountPayment,
		TypeOwnBankAccountIncoming,
		TypeOwnCounterPayment,
		TypeOwnCounterIncoming,
		TypeCorrection,
		TypeFeeReceived,
		TypeForeignCurrencyPayment,
		TypeCardFee,
		TypeDirectDebit,
		TypeDirectDebitIncoming,
		TypeDirectDebitOutgoing,
		TypeDirectDebitIncomingOtherBank,
		TypeInstantIncoming,
		TypeInstantOutgoing,
		TypeMortgageInsuranceFee,
		TypeInstantEuroIncoming,
		TypeInstantEuroOutgoing:
		return true
	default:
		return false
	}
}
