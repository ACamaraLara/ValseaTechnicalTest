package tester

import "time"

type Account struct {
	ID      string  `json:"id" bson:"_id"`
	Owner   string  `json:"owner" bson:"owner"`
	Balance float64 `json:"balance" bson:"balance"`
}

// Keyed by owner name to don't make the userstore the accountID to perform transactions
// It also updates in real time accounts information after transactions and transfers.
type AccountsRepo map[string]*Account

func NewAccountsRepo() AccountsRepo {
	return make(AccountsRepo)
}

func (acc *Account) Deposit(amount float64) {
	acc.Balance += amount
}

func (acc *Account) Withdraw(amount float64) {
	acc.Balance -= amount
}

type Transaction struct {
	ID        string    `json:"id" bson:"_id"`
	AccountID string    `json:"account_id" bson:"account_id"`
	Type      string    `json:"type" bson:"type"`
	Amount    float64   `json:"amount" bson:"amount"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
