package bank

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAccountInput(t *testing.T) {
	tests := []struct {
		owner          string
		initialBalance float64
		expectedError  error
	}{
		{"John Doe", 1000.0, nil},
		{"", 1000.0, ErrEmptyOwnerName},
		{"John Doe", -500.0, ErrNegativeInitialBalance},
	}

	for _, test := range tests {
		t.Run(test.owner, func(t *testing.T) {
			err := ValidateAccountInput(test.owner, test.initialBalance)
			if test.expectedError != nil {
				assert.True(t, errors.Is(err, test.expectedError))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTransaction(t *testing.T) {
	tests := []struct {
		txType        string
		amount        float64
		expectedError error
	}{
		{DepositTransactionType, 500.0, nil},
		{WithdrawalTransactionType, 500.0, nil},
		{"invalid", 500.0, ErrInvalidTransaction},
		{DepositTransactionType, 0.0, ErrZeroTransactionAmount},
		{WithdrawalTransactionType, 0.0, ErrZeroTransactionAmount},
	}

	for _, test := range tests {
		t.Run(test.txType, func(t *testing.T) {
			err := ValidateTransaction(test.txType, test.amount)
			if test.expectedError != nil {
				assert.True(t, errors.Is(err, test.expectedError))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTransfer(t *testing.T) {
	tests := []struct {
		fromAccountID string
		toAccountID   string
		amount        float64
		expectedError error
	}{
		{"1", "2", 500.0, nil},
		{"1", "1", 500.0, ErrSameSourceDestination},
		{"1", "2", -500.0, ErrZeroTransactionAmount},
	}

	for _, test := range tests {
		t.Run(test.fromAccountID+"->"+test.toAccountID, func(t *testing.T) {
			err := ValidateTransfer(test.fromAccountID, test.toAccountID, test.amount)
			if test.expectedError != nil {
				assert.True(t, errors.Is(err, test.expectedError))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
