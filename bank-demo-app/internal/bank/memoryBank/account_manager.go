package memoryBank

import (
	"bank-demo-app/internal/bank"

	"sync"

	"github.com/google/uuid"
)

type AccountManager struct {
	mu       sync.RWMutex            // Protect from race conditions.
	Accounts map[string]bank.Account // Keyed by AccountID. Using map because it's faster than an slice to get speciffic account.
}

func (am *AccountManager) CreateAccount(owner string, initialBalance float64) (*bank.Account, error) {
	// Assuming that there could be different accounts for a same owner name but with different IDs,
	// no need to check if account owner exists.
	account := bank.Account{
		ID:      uuid.New().String(),
		Owner:   owner,
		Balance: initialBalance,
	}

	am.mu.Lock()
	am.Accounts[account.ID] = account
	am.mu.Unlock()

	return &account, nil
}

func (am *AccountManager) GetAccountByID(accountID string) (*bank.Account, error) {
	am.mu.RLock()
	account, exists := am.Accounts[accountID]
	am.mu.RUnlock()

	if !exists {
		return &bank.Account{}, bank.AccountNotFoundError(accountID)
	}

	return &account, nil
}

func (am *AccountManager) ListAccounts() []bank.Account {
	am.mu.RLock()
	accounts := make([]bank.Account, 0, len(am.Accounts))
	for _, account := range am.Accounts {
		accounts = append(accounts, account)
	}
	am.mu.RUnlock()

	return accounts
}

func (am *AccountManager) TransferBetweenAccounts(fromAccountID, toAccountID string, amount float64) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	fromAccount, exists := am.Accounts[fromAccountID]
	if !exists {
		return bank.TransferSourceNotFoundError(fromAccountID)
	}

	toAccount, exists := am.Accounts[toAccountID]
	if !exists {
		return bank.TransferDestinationNotFoundError(toAccountID)
	}

	// Withdraw from one account and deposit into the other one.
	if err := fromAccount.Withdraw(amount); err != nil {
		return err
	}
	toAccount.Deposit(amount)

	am.Accounts[fromAccountID] = fromAccount
	am.Accounts[toAccountID] = toAccount

	return nil
}

func (am *AccountManager) PerformTransaction(account *bank.Account, transactionType string, amount float64) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if err := account.UpdateBalance(transactionType, amount); err != nil {
		return err
	}
	am.Accounts[account.ID] = *account
	return nil
}
