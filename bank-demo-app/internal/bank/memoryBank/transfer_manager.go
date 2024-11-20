package memoryBank

import (
	"bank-demo-app/internal/bank"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TransactionManager struct {
	mu           sync.RWMutex                  // Protect against race conditions
	Transactions map[string][]bank.Transaction // Keyed by AccountID
}

func (tm *TransactionManager) CreateTransaction(accountID, transactionType string, amount float64) (*bank.Transaction, error) {
	// Create the transaction
	transaction := bank.Transaction{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Type:      transactionType,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	// Store the transaction
	tm.mu.Lock()
	tm.Transactions[accountID] = append(tm.Transactions[accountID], transaction)
	tm.mu.Unlock()

	return &transaction, nil
}

func (tm *TransactionManager) GetTransactionsByAccountID(accountID string) ([]bank.Transaction, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	transactions, exists := tm.Transactions[accountID]
	if !exists {
		return nil, bank.NoTransactionsForAccountError(accountID)
	}

	return transactions, nil
}
