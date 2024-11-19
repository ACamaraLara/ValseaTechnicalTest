package bank

import (
	"errors"
	"fmt"
)

var (
	// Account errors.
	ErrAccountNotFound              = errors.New("account not found")
	ErrEmptyOwnerName               = errors.New("owner name cannot be empty")
	ErrSameSourceDestinationAccount = errors.New("source and destination account cannot be the same")
	ErrNoTransactionsForAccount     = errors.New("no transactions found in provided account")

	// Transaction errors.
	ErrTransactionNotFound     = errors.New("transaction not found")
	ErrInvalidTransaction      = errors.New("invalid transaction type")
	ErrTransactionTypeRequired = errors.New("transaction type is required")
	ErrZeroTransactionAmount   = errors.New("transaction amount must be greater than zero")

	// Balance errors
	ErrNegativeAmount         = errors.New("amount must be greater than zero")
	ErrNegativeInitialBalance = errors.New("initial balance cannot be negative")
	ErrInsufficientFunds      = errors.New("insufficient funds")

	// Transfer errors
	ErrTransferSourceNotFound      = errors.New("transfer account source not found")
	ErrTransferDestinationNotFound = errors.New("transfer account source not found")
)

// Helper functions for error wrapping.

// Account errors.
func AccountNotFoundError(accountID string) error {
	return fmt.Errorf("%w: account ID %s", ErrAccountNotFound, accountID)
}

func EmptyOwnerNameError() error {
	return fmt.Errorf("%w", ErrEmptyOwnerName)
}

func SameSourceDestinationAccountError(accountID string) error {
	return fmt.Errorf("%w: account ID %s", ErrSameSourceDestinationAccount, accountID)
}

func NoTransactionsForAccountError(accountID string) error {
	return fmt.Errorf("%w: account ID %s", ErrNoTransactionsForAccount, accountID)
}

// Transaction errors.
func TransactionNotFoundError(transactionID string) error {
	return fmt.Errorf("%w: transaction ID %s", ErrTransactionNotFound, transactionID)
}

func InvalidTransactionError(transactionType string) error {
	return fmt.Errorf("%w: transaction type %s", ErrInvalidTransaction, transactionType)
}

func TransactionTypeRequiredError() error {
	return fmt.Errorf("%w", ErrTransactionTypeRequired)
}

func ZeroTransactionAmountError() error {
	return fmt.Errorf("%w", ErrZeroTransactionAmount)
}

// Balance errors.
func NegativeAmountError(amount float64) error {
	return fmt.Errorf("%w: amount %.2f", ErrNegativeAmount, amount)
}

func NegativeInitialBalanceError(initialBalance float64) error {
	return fmt.Errorf("%w: initial balance %.2f", ErrNegativeInitialBalance, initialBalance)
}

func InsufficientFundsError(accountID string, balance, amount float64) error {
	return fmt.Errorf("%w: account ID %s, balance %.2f, attempted %.2f", ErrInsufficientFunds, accountID, balance, amount)
}

// Transfer errors.
func TransferSourceNotFoundError(accountID string) error {
	return fmt.Errorf("%w: account ID %s", ErrTransferSourceNotFound, accountID)
}

func TransferDestinationNotFoundError(accountID string) error {
	return fmt.Errorf("%w: account ID %s", ErrTransferDestinationNotFound, accountID)
}
