package restServer

import (
	"bank-demo-app/internal/bank"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// BankStore defines the methods required for managing accounts and transactions.
type BankStore interface {
	// Account operations
	CreateAccount(owner string, initialBalance float64) (*bank.Account, error)
	GetAccountByID(id string) (*bank.Account, error)
	ListAccounts() []bank.Account

	// Transaction operations
	PerformTransaction(accountID string, txType string, amount float64) (*bank.Transaction, error)
	GetTransactionsByAccountID(accountID string) ([]bank.Transaction, error)

	// Transfer operations
	TransferFunds(fromAccountID, toAccountID string, amount float64) error
}

// This is the default status handler that will be used to check if the REST server is up.
func statusHandler(c *gin.Context) {
	log.Info().Msg("Called GET status method.")
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Expected GET method!",
		})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func createAccountHandler(bankStore BankStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Owner          string  `json:"owner"`
			InitialBalance float64 `json:"initial_balance"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Error().Err(err).Msg("Invalid request body while creating account")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		log.Info().Str("owner", request.Owner).Float64("initial_balance", request.InitialBalance).Msg("Creating account")

		account, err := bankStore.CreateAccount(request.Owner, request.InitialBalance)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create account")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Info().Str("account_id", account.ID).Str("owner", account.Owner).Float64("initial_balance", account.Balance).Msg("Account created successfully")
		c.JSON(http.StatusCreated, account)
	}
}

func listAccountsHandler(bankStore BankStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info().Msg("Listing all accounts")

		accounts := bankStore.ListAccounts()

		log.Info().Int("accounts_count", len(accounts)).Msg("Accounts listed successfully")
		c.JSON(http.StatusOK, accounts)
	}
}

func getAccountByIDHandler(bankStore BankStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Param("id")

		log.Info().Str("account_id", accountID).Msg("Retrieving account details")

		account, err := bankStore.GetAccountByID(accountID)
		if err != nil {
			log.Warn().Err(err).Str("account_id", accountID).Msg("Account not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}

		log.Info().Str("account_id", accountID).Str("owner", account.Owner).Float64("balance", account.Balance).Msg("Account retrieved successfully")
		c.JSON(http.StatusOK, account)
	}
}

func performTransactionHandler(bankStore BankStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Param("id")
		var request struct {
			Type   string  `json:"type"`
			Amount float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Warn().Err(err).Str("account_id", accountID).Msg("Invalid request body for transaction")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if request.Type != "deposit" && request.Type != "withdrawal" {
			log.Warn().Str("account_id", accountID).Msg("Invalid transaction type")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
			return
		}

		log.Info().Str("account_id", accountID).Str("transaction_type", request.Type).Float64("amount", request.Amount).Msg("Creating transaction")

		transaction, err := bankStore.PerformTransaction(accountID, request.Type, request.Amount)
		if err != nil {
			log.Error().Err(err).Str("account_id", accountID).Msg("Failed to create transaction")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Info().Str("account_id", accountID).Str("transaction_id", transaction.ID).Float64("amount", transaction.Amount).Msg("Transaction created successfully")
		c.JSON(http.StatusCreated, transaction)
	}
}

func getTransactionsByAccountIDHandler(bankStore BankStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Param("id")

		log.Info().Str("account_id", accountID).Msg("Retrieving transactions")

		transactions, err := bankStore.GetTransactionsByAccountID(accountID)
		if err != nil {
			log.Warn().Err(err).Str("account_id", accountID).Msg("Transactions not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Transactions not found"})
			return
		}

		log.Info().Str("account_id", accountID).Int("transaction_count", len(transactions)).Msg("Transactions retrieved successfully")
		c.JSON(http.StatusOK, transactions)
	}
}

func transferFundsHandler(bankStore BankStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			FromAccountID string  `json:"from_account_id"`
			ToAccountID   string  `json:"to_account_id"`
			Amount        float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Warn().Err(err).Msg("Invalid request body for transfer")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		log.Info().Str("from_account_id", request.FromAccountID).Str("to_account_id", request.ToAccountID).Float64("amount", request.Amount).Msg("Initiating fund transfer")

		err := bankStore.TransferFunds(request.FromAccountID, request.ToAccountID, request.Amount)
		if err != nil {
			log.Error().Err(err).Str("from_account_id", request.FromAccountID).Str("to_account_id", request.ToAccountID).Msg("Failed to transfer funds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Info().Str("from_account_id", request.FromAccountID).Str("to_account_id", request.ToAccountID).Float64("amount", request.Amount).Msg("Transfer successful")
		c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
	}
}
