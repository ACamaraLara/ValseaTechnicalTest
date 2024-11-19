package bank

import (
	"time"
)

const (
	DepositTransactionType    = "deposit"
	WithdrawalTransactionType = "withdrawal"
)

// Account represents a bank account owner information.
type Account struct {
	ID      string  `json:"id" bson:"_id"`
	Owner   string  `json:"owner" bson:"owner"`
	Balance float64 `json:"balance" bson:"balance"`
}

func (acc *Account) Deposit(amount float64) {
	acc.Balance += amount
}

func (acc *Account) Withdraw(amount float64) error {
	if acc.Balance < amount {
		return InsufficientFundsError(acc.ID, acc.Balance, amount)
	}
	acc.Balance -= amount
	return nil
}

// Transaction represents a financial transaction associated with an account.
type Transaction struct {
	ID        string    `json:"id" bson:"_id"`
	AccountID string    `json:"account_id" bson:"account_id"`
	Type      string    `json:"type" bson:"type"`
	Amount    float64   `json:"amount" bson:"amount"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
