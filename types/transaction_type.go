package types

import "database/sql/driver"

const (
	typeIncomingInsideBank               = "Příjem převodem uvnitř banky"
	typeOutgoingInsideBank               = "Platba převodem uvnitř banky"
	typeCashDepositAtCounter             = "Vklad pokladnou"
	typeCashWithdrawalAtCounter          = "Výběr pokladnou"
	typeCashDeposit                      = "Vklad v hotovosti"
	typeCashWithdrawal                   = "Výběr v hotovosti"
	typePayment                          = "Platba"
	typeIncoming                         = "Příjem"
	typeCashlessPayment                  = "Bezhotovostní platba"
	typeCashlessIncoming                 = "Bezhotovostní příjem"
	typeCardPayment                      = "Platba kartou"
	typeLoanInterest                     = "Úrok z úvěru"
	typePenaltyFee                       = "Sankční poplatek"
	typeCourierOutgoing                  = "Posel – předání"
	typeCourierIncoming                  = "Posel – příjem"
	typeTransferInsideAccount            = "Převod uvnitř konta"
	typeInterestCredited                 = "Připsaný úrok"
	typeInterestPaid                     = "Vyplacený úrok"
	typeInterestTax                      = "Odvod daně z úroků"
	typeRecordedInterest                 = "Evidovaný úrok"
	typeFee                              = "Poplatek"
	typeRecordedFee                      = "Evidovaný poplatek"
	typeInterBankAccountTransferOutgoing = "Převod mezi bankovními konty (platba)"
	typeInterBankAccountTransferIncoming = "Převod mezi bankovními konty (příjem)"
	typeUnidentifiedBankAccountPayment   = "Neidentifikovaná platba z bankovního konta"
	typeUnidentifiedBankAccountIncoming  = "Neidentifikovaný příjem na bankovní konto"
	typeOwnBankAccountPayment            = "Vlastní platba z bankovního konta"
	typeOwnBankAccountIncoming           = "Vlastní příjem na bankovní konto"
	typeOwnCounterPayment                = "Vlastní platba pokladnou"
	typeOwnCounterIncoming               = "Vlastní příjem pokladnou"
	typeCorrection                       = "Opravný pohyb"
	typeFeeReceived                      = "Přijatý poplatek"
	typeForeignCurrencyPayment           = "Platba v jiné měně"
	typeCardFee                          = "Poplatek – platební karta"
	typeDirectDebit                      = "Inkaso"
	typeDirectDebitIncoming              = "Inkaso ve prospěch účtu"
	typeDirectDebitOutgoing              = "Inkaso z účtu"
	typeDirectDebitIncomingOtherBank     = "Příjem inkasa z cizí banky"
	typeInstantIncoming                  = "Okamžitá příchozí platba"
	typeInstantOutgoing                  = "Okamžitá odchozí platba"
	typeMortgageInsuranceFee             = "Poplatek - pojištění hypotéky"
	typeInstantEuroIncoming              = "Okamžitá příchozí Europlatba"
	typeInstantEuroOutgoing              = "Okamžitá odchozí Europlatba"
)

// TransactionType is the localized transaction category reported by Fio
// bank. Unknown values are preserved and can be detected with [TransactionType.IsKnown].
type TransactionType struct {
	innerType string
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

// Value implements [driver.Valuer] by returning the transaction type's
// original string representation.
func (receiver TransactionType) Value() (driver.Value, error) {
	return receiver.String(), nil
}

// IsIncoming reports whether the type represents incoming funds.
func (receiver TransactionType) IsIncoming() bool {
	switch receiver.innerType {
	case
		typeIncomingInsideBank,
		typeIncoming,
		typeCashlessIncoming,
		typeCourierIncoming,
		typeInterBankAccountTransferIncoming,
		typeUnidentifiedBankAccountIncoming,
		typeOwnBankAccountIncoming,
		typeOwnCounterIncoming,
		typeFeeReceived,
		typeDirectDebitIncoming,
		typeDirectDebitIncomingOtherBank,
		typeInstantIncoming,
		typeInstantEuroIncoming:
		return true
	default:
		return false
	}
}

// IsOutgoing reports whether the type represents outgoing funds.
func (receiver TransactionType) IsOutgoing() bool {
	switch receiver.innerType {
	case
		typeOutgoingInsideBank,
		typeCashWithdrawalAtCounter,
		typeCashWithdrawal,
		typePayment,
		typeCashlessPayment,
		typeCardPayment,
		typeLoanInterest,
		typePenaltyFee,
		typeCourierOutgoing,
		typeInterBankAccountTransferOutgoing,
		typeUnidentifiedBankAccountPayment,
		typeOwnBankAccountPayment,
		typeOwnCounterPayment,
		typeFee,
		typeCardFee,
		typeDirectDebitOutgoing,
		typeMortgageInsuranceFee,
		typeForeignCurrencyPayment,
		typeInstantOutgoing,
		typeInstantEuroOutgoing:
		return true
	default:
		return false
	}
}

// IsImmediate reports whether the type represents an instant payment.
func (receiver TransactionType) IsImmediate() bool {
	switch receiver.innerType {
	case
		typeInstantIncoming,
		typeInstantOutgoing,
		typeInstantEuroIncoming,
		typeInstantEuroOutgoing:
		return true
	default:
		return false
	}
}

// IsInternal reports whether the type represents a transfer within Fio bank
// or between the account owner's accounts.
func (receiver TransactionType) IsInternal() bool {
	switch receiver.innerType {
	case
		typeIncomingInsideBank,
		typeOutgoingInsideBank,
		typeTransferInsideAccount,
		typeOwnBankAccountPayment,
		typeOwnBankAccountIncoming:
		return true
	default:
		return false
	}
}

// IsCardPayment reports whether the type represents a card payment.
func (receiver TransactionType) IsCardPayment() bool {
	return receiver.innerType == typeCardPayment
}

// IsCash reports whether the type represents a cash transaction.
func (receiver TransactionType) IsCash() bool {
	switch receiver.innerType {
	case
		typeCashDepositAtCounter,
		typeCashWithdrawalAtCounter,
		typeCashDeposit,
		typeCashWithdrawal,
		typeOwnCounterPayment,
		typeOwnCounterIncoming:
		return true
	default:
		return false
	}
}

// IsDirectDebit reports whether the type represents a direct debit.
func (receiver TransactionType) IsDirectDebit() bool {
	switch receiver.innerType {
	case
		typeDirectDebit,
		typeDirectDebitIncoming,
		typeDirectDebitOutgoing,
		typeDirectDebitIncomingOtherBank:
		return true
	default:
		return false
	}
}

// IsFee reports whether the type represents a fee.
func (receiver TransactionType) IsFee() bool {
	switch receiver.innerType {
	case
		typePenaltyFee,
		typeFee,
		typeRecordedFee,
		typeFeeReceived,
		typeCardFee,
		typeMortgageInsuranceFee:
		return true
	default:
		return false
	}
}

// IsInterest reports whether the type represents interest or interest tax.
func (receiver TransactionType) IsInterest() bool {
	switch receiver.innerType {
	case
		typeLoanInterest,
		typeInterestCredited,
		typeInterestPaid,
		typeInterestTax,
		typeRecordedInterest:
		return true
	default:
		return false
	}
}

// IsCorrection reports whether the type represents a corrective entry.
func (receiver TransactionType) IsCorrection() bool {
	return receiver.innerType == typeCorrection
}

// IsKnown reports whether the type is recognized by this version of the
// package.
func (receiver TransactionType) IsKnown() bool {
	switch receiver.innerType {
	case
		typeIncomingInsideBank,
		typeOutgoingInsideBank,
		typeCashDepositAtCounter,
		typeCashWithdrawalAtCounter,
		typeCashDeposit,
		typeCashWithdrawal,
		typePayment,
		typeIncoming,
		typeCashlessPayment,
		typeCashlessIncoming,
		typeCardPayment,
		typeLoanInterest,
		typePenaltyFee,
		typeCourierOutgoing,
		typeCourierIncoming,
		typeTransferInsideAccount,
		typeInterestCredited,
		typeInterestPaid,
		typeInterestTax,
		typeRecordedInterest,
		typeFee,
		typeRecordedFee,
		typeInterBankAccountTransferOutgoing,
		typeInterBankAccountTransferIncoming,
		typeUnidentifiedBankAccountPayment,
		typeUnidentifiedBankAccountIncoming,
		typeOwnBankAccountPayment,
		typeOwnBankAccountIncoming,
		typeOwnCounterPayment,
		typeOwnCounterIncoming,
		typeCorrection,
		typeFeeReceived,
		typeForeignCurrencyPayment,
		typeCardFee,
		typeDirectDebit,
		typeDirectDebitIncoming,
		typeDirectDebitOutgoing,
		typeDirectDebitIncomingOtherBank,
		typeInstantIncoming,
		typeInstantOutgoing,
		typeMortgageInsuranceFee,
		typeInstantEuroIncoming,
		typeInstantEuroOutgoing:
		return true
	default:
		return false
	}
}
