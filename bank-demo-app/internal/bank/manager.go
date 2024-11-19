package bank

import "time"

// Account represents a bank account owner information.
type Account struct {
	ID      string  `json:"id" bson:"_id"`
	Owner   string  `json:"owner" bson:"owner"`
	Balance float64 `json:"balance" bson:"balance"`
}

// Transaction represents a financial transaction associated with an account.
type Transaction struct {
	ID        string    `json:"id" bson:"_id"`
	AccountID string    `json:"account_id" bson:"account_id"`
	Type      string    `json:"type" bson:"type"`
	Amount    float64   `json:"amount" bson:"amount"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// BankStore defines the methods required for managing accounts and transactions.
type BankStore interface {
	// Account operations
	CreateAccount(owner string, initialBalance float64) (*Account, error)
	GetAccountByID(id string) (*Account, error)
	ListAccounts() ([]Account, error)

	// Transaction operations
	CreateTransaction(accountID string, txType string, amount float64) (*Transaction, error)
	GetTransactionsByAccountID(accountID string) ([]Transaction, error)

	// Transfer operations
	TransferFunds(fromAccountID, toAccountID string, amount float64) error
}

type BankManager struct {
	Store BankStore
}

func NewBankManager(store BankStore) *BankManager {
	return &BankManager{Store: store}
}
