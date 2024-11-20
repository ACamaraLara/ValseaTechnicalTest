package bank

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	account := &Account{
		ID:      "1",
		Owner:   "John Doe",
		Balance: 1000.0,
	}

	depositAmount := 500.0
	account.Deposit(depositAmount)

	assert.Equal(t, 1500.0, account.Balance)
}

func TestWithdraw(t *testing.T) {
	account := &Account{
		ID:      "1",
		Owner:   "John Doe",
		Balance: 1000.0,
	}

	withdrawAmount := 500.0
	err := account.Withdraw(withdrawAmount)
	assert.NoError(t, err)

	assert.Equal(t, 500.0, account.Balance)

	invalidWithdrawAmount := 600.0
	err = account.Withdraw(invalidWithdrawAmount)

	assert.True(t, errors.Is(err, ErrInsufficientFunds))
}

func TestUpdateBalance_Deposit(t *testing.T) {
	account := &Account{
		ID:      "1",
		Owner:   "John Doe",
		Balance: 1000.0,
	}

	txType := "deposit"
	amount := 500.0
	err := account.UpdateBalance(txType, amount)
	assert.NoError(t, err)

	assert.Equal(t, 1500.0, account.Balance)
}

func TestUpdateBalance_Withdrawal(t *testing.T) {
	account := &Account{
		ID:      "1",
		Owner:   "John Doe",
		Balance: 1000.0,
	}

	txType := "withdrawal"
	amount := 500.0
	err := account.UpdateBalance(txType, amount)
	assert.NoError(t, err)

	assert.Equal(t, 500.0, account.Balance)

	invalidAmount := 600.0
	err = account.UpdateBalance(txType, invalidAmount)
	assert.True(t, errors.Is(err, ErrInsufficientFunds))
}

func TestUpdateBalance_InvalidTransactionType(t *testing.T) {
	account := &Account{
		ID:      "1",
		Owner:   "John Doe",
		Balance: 1000.0,
	}

	txType := "invalid"
	amount := 500.0
	err := account.UpdateBalance(txType, amount)

	assert.True(t, errors.Is(err, ErrInvalidTransaction))
}
