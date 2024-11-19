package memoryBank

import "bank-demo-app/internal/bank"

type BankStore struct {
	accManager      *AccountManager
	transactManager *TransactionManager
}

func NewBankStore() *BankStore {
	// Initialize the AccountManager and TransactionManager
	accManager := &AccountManager{
		Accounts: make(map[string]bank.Account), // Initialize the map to hold accounts
	}

	transactManager := &TransactionManager{
		Transactions: make(map[string][]bank.Transaction), // Initialize the map to hold transactions
	}

	// Create and return the BankStore with both managers
	return &BankStore{
		accManager:      accManager,
		transactManager: transactManager,
	}
}

func (bs *BankStore) CreateAccount(owner string, initialBalance float64) (*bank.Account, error) {
	return bs.accManager.CreateAccount(owner, initialBalance)
}

func (bs *BankStore) GetAccountByID(id string) (*bank.Account, error) {
	return bs.accManager.GetAccountByID(id)
}

func (bs *BankStore) ListAccounts() []bank.Account {
	return bs.accManager.ListAccounts()
}

func (bs *BankStore) PerformTransaction(accountID string, txType string, amount float64) (*bank.Transaction, error) {
	account, err := bs.accManager.GetAccountByID(accountID)
	if err != nil {
		return &bank.Transaction{}, err
	}

	if err := bs.accManager.PerformTransaction(account, txType, amount); err != nil {
		return &bank.Transaction{}, err
	}

	return bs.transactManager.CreateTransaction(accountID, txType, amount)
}

func (bs *BankStore) GetTransactionsByAccountID(accountID string) ([]bank.Transaction, error) {
	return bs.transactManager.GetTransactionsByAccountID(accountID)
}

func (bs *BankStore) TransferFunds(fromAccountID, toAccountID string, amount float64) error {
	return bs.accManager.TransferBetweenAccounts(fromAccountID, toAccountID, amount)
}
