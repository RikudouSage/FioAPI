package types

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

type TransactionType struct {
	innerType string
}

func (receiver TransactionType) MarshalText() ([]byte, error) {
	return []byte(receiver.innerType), nil
}

func (receiver *TransactionType) UnmarshalText(text []byte) error {
	*receiver = TransactionType{innerType: string(text)}
	return nil
}

func (receiver TransactionType) String() string {
	return receiver.innerType
}

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

func (receiver TransactionType) IsCardPayment() bool {
	return receiver.innerType == typeCardPayment
}

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

func (receiver TransactionType) IsCorrection() bool {
	return receiver.innerType == typeCorrection
}

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
