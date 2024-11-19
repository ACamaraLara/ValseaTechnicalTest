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
	if owner == "" {
		return &bank.Account{}, bank.ErrEmptyOwnerName
	}
	if initialBalance < 0 {
		return &bank.Account{}, bank.ErrNegativeInitialBalance
	}

	// Assuming that there could be different accounts for a same owner name but with different IDs,
	// no need to check if account owner exists.
	accountID := uuid.New().String()
	account := bank.Account{
		ID:      accountID,
		Owner:   owner,
		Balance: initialBalance,
	}

	am.mu.Lock()
	am.Accounts[accountID] = account
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
	if amount <= 0 {
		return bank.ErrZeroTransactionAmount
	}

	if fromAccountID == toAccountID {
		return bank.ErrSameSourceDestination
	}

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

	if fromAccount.Balance < amount {
		return bank.InsufficientFundsError(fromAccountID, fromAccount.Balance, amount)
	}

	fromAccount.Withdraw(amount)
	toAccount.Deposit(amount)

	am.Accounts[fromAccountID] = fromAccount
	am.Accounts[toAccountID] = toAccount

	return nil
}

func (am *AccountManager) PerformTransaction(account *bank.Account, transactionType string, amount float64) error {
	if amount <= 0 {
		return bank.ErrZeroTransactionAmount
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	if transactionType == bank.DepositTransactionType {
		account.Deposit(amount)
	} else if transactionType == bank.WithdrawalTransactionType {
		if err := account.Withdraw(amount); err != nil {
			return err
		}
	} else {
		return bank.InvalidTransactionError(transactionType)
	}
	am.Accounts[account.ID] = *account
	return nil
}
