package bank

func ValidateAccountInput(owner string, initialBalance float64) error {
	if owner == "" {
		return ErrEmptyOwnerName
	}
	if initialBalance < 0 {
		return ErrNegativeInitialBalance
	}
	return nil
}

func ValidateTransaction(txType string, amount float64) error {
	if txType != DepositTransactionType && txType != WithdrawalTransactionType {
		return InvalidTransactionError(txType)
	}
	if amount <= 0 {
		return ErrZeroTransactionAmount
	}
	return nil
}

func ValidateTransfer(fromAccountID, toAccountID string, amount float64) error {
	if amount <= 0 {
		return ErrZeroTransactionAmount
	}
	if fromAccountID == toAccountID {
		return SameSourceDestinationAccountError(fromAccountID)
	}
	return nil
}
