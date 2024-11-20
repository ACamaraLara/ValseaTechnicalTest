package memoryBank

import (
	"bank-demo-app/internal/bank"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	bankStore := NewBankStore()

	owner := "John Doe"
	initialBalance := 1000.0
	account, err := bankStore.CreateAccount(owner, initialBalance)

	assert.NoError(t, err)

	assert.NotEmpty(t, account.ID)
	assert.Equal(t, owner, account.Owner)
	assert.Equal(t, initialBalance, account.Balance)
}

func TestGetAccountByID(t *testing.T) {
	bankStore := NewBankStore()

	owner := "John Doe"
	initialBalance := 1000.0
	account, err := bankStore.CreateAccount(owner, initialBalance)
	assert.NoError(t, err)

	retrievedAccount, err := bankStore.GetAccountByID(account.ID)
	assert.NoError(t, err)
	assert.Equal(t, account, retrievedAccount)

	nonExistentID := uuid.New().String()
	retrievedAccount, err = bankStore.GetAccountByID(nonExistentID)
	assert.Error(t, err)
	assert.Equal(t, bank.Account{}, *retrievedAccount)
}

func TestListAccounts(t *testing.T) {
	bankStore := NewBankStore()

	// Create two accounts
	owner1 := "Alex Camara"
	initialBalance1 := 1000.0
	account1, err := bankStore.CreateAccount(owner1, initialBalance1)
	assert.NoError(t, err)

	owner2 := "Donald Trump"
	initialBalance2 := 500.0
	account2, err := bankStore.CreateAccount(owner2, initialBalance2)
	assert.NoError(t, err)

	accounts := bankStore.ListAccounts()

	assert.Len(t, accounts, 2)
	assert.Contains(t, accounts, *account1)
	assert.Contains(t, accounts, *account2)
}

func TestPerformTransaction(t *testing.T) {
	bankStore := NewBankStore()

	owner := "Alex Camara"
	initialBalance := 1000.0
	account, err := bankStore.CreateAccount(owner, initialBalance)
	assert.NoError(t, err)

	txType := "deposit"
	amount := 500.0
	transaction, err := bankStore.PerformTransaction(account.ID, txType, amount)
	assert.NoError(t, err)

	assert.Equal(t, account.ID, transaction.AccountID)
	assert.Equal(t, txType, transaction.Type)
	assert.Equal(t, amount, transaction.Amount)

	accountAfter, err := bankStore.GetAccountByID(account.ID)
	assert.NoError(t, err)
	assert.Equal(t, initialBalance+amount, accountAfter.Balance)

	txType = "withdrawal"
	amount = 200.0
	_, err = bankStore.PerformTransaction(account.ID, txType, amount)
	assert.NoError(t, err)

	accountAfter, err = bankStore.GetAccountByID(account.ID)
	assert.NoError(t, err)
	assert.Equal(t, initialBalance+300.0, accountAfter.Balance)

	invalidTxType := "invalid"
	transaction, err = bankStore.PerformTransaction(account.ID, invalidTxType, amount)
	assert.Error(t, err)
	assert.Equal(t, bank.Transaction{}, *transaction)
}

func TestGetTransactionsByAccountID(t *testing.T) {
	bankStore := NewBankStore()

	owner := "Alex Camara"
	initialBalance := 1000.0
	account, err := bankStore.CreateAccount(owner, initialBalance)
	assert.NoError(t, err)

	txType := "deposit"
	amount := 500.0
	_, err = bankStore.PerformTransaction(account.ID, txType, amount)
	assert.NoError(t, err)

	transactions, err := bankStore.GetTransactionsByAccountID(account.ID)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)

	assert.Equal(t, account.ID, transactions[0].AccountID)
	assert.Equal(t, txType, transactions[0].Type)
	assert.Equal(t, amount, transactions[0].Amount)

	nonExistentID := uuid.New().String()
	transactions, err = bankStore.GetTransactionsByAccountID(nonExistentID)
	assert.Error(t, err)
	assert.Empty(t, transactions)
}

func TestTransferFunds(t *testing.T) {
	bankStore := NewBankStore()

	// Create two accounts for the transfer test
	owner1 := "John Doe"
	initialBalance1 := 1000.0
	account1, err := bankStore.CreateAccount(owner1, initialBalance1)
	assert.NoError(t, err)

	owner2 := "Jane Doe"
	initialBalance2 := 1500.0
	account2, err := bankStore.CreateAccount(owner2, initialBalance2)
	assert.NoError(t, err)

	// Perform a transfer of 200 from account1 to account2
	transferAmount := 200.0
	err = bankStore.TransferFunds(account1.ID, account2.ID, transferAmount)
	assert.NoError(t, err)

	// Assert the balances after the transfer
	account1After, err := bankStore.GetAccountByID(account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, initialBalance1-transferAmount, account1After.Balance)

	account2After, err := bankStore.GetAccountByID(account2.ID)
	assert.NoError(t, err)
	assert.Equal(t, initialBalance2+transferAmount, account2After.Balance)

	// Test transfer with invalid account ID
	invalidID := uuid.New().String()
	err = bankStore.TransferFunds(account1.ID, invalidID, transferAmount)
	assert.Error(t, err)

	// Test transfer with insufficient funds
	insufficientBalanceAmount := initialBalance1 + 500.0
	err = bankStore.TransferFunds(account1.ID, account2.ID, insufficientBalanceAmount)
	assert.Error(t, err)
}
