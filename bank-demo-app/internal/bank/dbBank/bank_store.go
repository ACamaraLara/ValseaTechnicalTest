package dbBank

import (
	"bank-demo-app/internal/bank"
	"bank-demo-app/internal/mongodb"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	accountsCollection     string = "accounts"
	transactionsCollection string = "transactions"
)

type BankStore struct {
	dbClient *mongodb.MongoDBClient
}

func NewBankStore(ctx context.Context, dbConf *mongodb.MongoConfig) *BankStore {
	mongoClient := mongodb.NewMongoDBClient(dbConf)

	if err := mongoClient.ConnectMongoClient(ctx); err != nil {
		return nil
	}
	// Workarround for testing application, this would be refactored and way more clean.
	mongoClient.GetCollections([]string{accountsCollection, transactionsCollection})

	bankStore := &BankStore{dbClient: mongoClient}

	return bankStore
}

func (bs *BankStore) CreateAccount(owner string, initialBalance float64) (*bank.Account, error) {
	accountID := uuid.New().String()
	account := bank.Account{
		ID:      accountID,
		Owner:   owner,
		Balance: initialBalance,
	}

	collection := bs.dbClient.Collections[accountsCollection]
	_, err := collection.InsertOne(context.Background(), account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return &account, nil
}

func (bs *BankStore) GetAccountByID(id string) (*bank.Account, error) {
	collection := bs.dbClient.Collections[accountsCollection]

	var account bank.Account
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&account)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, bank.AccountNotFoundError(id)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

func (bs *BankStore) ListAccounts() []bank.Account {
	collection := bs.dbClient.Collections[accountsCollection]

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil
	}
	defer cursor.Close(context.Background())

	// Decode all documents into a slice of bank.Account.
	var accounts []bank.Account
	if err := cursor.All(context.Background(), &accounts); err != nil {
		return nil
	}

	return accounts
}

func (bs *BankStore) PerformTransaction(accountID string, txType string, amount float64) (*bank.Transaction, error) {
	// Fetch the account
	account, err := bs.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	// Update the account balance
	if err := bs.updateAccountBalance(account, txType, amount); err != nil {
		return nil, err
	}

	// Create the transaction record
	transaction, err := bs.createTransaction(accountID, txType, amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil

}

// updateAccountBalance updates the account balance based on the transaction type
func (bs *BankStore) updateAccountBalance(account *bank.Account, txType string, amount float64) error {
	accountsCollection := bs.dbClient.Collections[accountsCollection]

	if err := account.UpdateBalance(txType, amount); err != nil {
		return err
	}

	_, err := accountsCollection.UpdateOne(context.Background(), bson.M{"_id": account.ID}, bson.M{"$set": bson.M{"balance": account.Balance}})
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	return nil
}

// createTransaction creates a new transaction record and stores it in the database
func (bs *BankStore) createTransaction(accountID, txType string, amount float64) (*bank.Transaction, error) {
	transactionsCollection := bs.dbClient.Collections[transactionsCollection]

	transaction := bank.Transaction{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Type:      txType,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	_, err := transactionsCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to insert transaction: %w", err)
	}

	return &transaction, nil
}

func (bs *BankStore) GetTransactionsByAccountID(accountID string) ([]bank.Transaction, error) {
	collection := bs.dbClient.Collections[transactionsCollection]

	// Find transactions by the accountID
	cursor, err := collection.Find(context.Background(), bson.M{"account_id": accountID})
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions: %w", err)
	}
	defer cursor.Close(context.Background())

	var transactions []bank.Transaction
	if err := cursor.All(context.Background(), &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode transactions: %w", err)
	}

	// If no transactions were found, return an error
	if len(transactions) == 0 {
		return nil, bank.NoTransactionsForAccountError(accountID)
	}

	return transactions, nil
}

func (bs *BankStore) TransferFunds(fromAccountID, toAccountID string, amount float64) error {
	// Get both accounts
	fromAccount, err := bs.GetAccountByID(fromAccountID)
	if err != nil {
		return err
	}

	toAccount, err := bs.GetAccountByID(toAccountID)
	if err != nil {
		return err
	}

	if err := bs.performTransfer(fromAccount, toAccount, amount); err != nil {
		return err
	}

	return nil
}

func (bs *BankStore) performTransfer(fromAccount, toAccount *bank.Account, amount float64) error {
	// Withdraw from the source account and deposit to the destination account
	if err := fromAccount.Withdraw(amount); err != nil {
		return err
	}
	toAccount.Deposit(amount)

	// Save the updated accounts back to the database.
	if err := bs.saveAccount(fromAccount); err != nil {
		return err
	}
	if err := bs.saveAccount(toAccount); err != nil {
		return err
	}

	return nil
}

func (bs *BankStore) saveAccount(account *bank.Account) error {
	collection := bs.dbClient.Collections[accountsCollection]
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": account.ID}, bson.M{"$set": account})
	return err
}
