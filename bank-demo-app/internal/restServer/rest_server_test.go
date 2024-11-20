package restServer

import (
	"bank-demo-app/internal/bank"
	"bank-demo-app/internal/bank/memoryBank"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BankRestAPITestSuite struct {
	suite.Suite
	router    *gin.Engine
	bankStore *memoryBank.BankStore
}

func (suite *BankRestAPITestSuite) SetupTest() {
	suite.bankStore = memoryBank.NewBankStore()
	routes := InitRestRoutes(suite.bankStore)
	suite.router = NewRouter(routes)
}

func (suite *BankRestAPITestSuite) TestStatusHandler() {
	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BankRestAPITestSuite) TestCreateAccountHandler() {
	account := createAccountRequest{
		Owner:          "Alex Camara",
		InitialBalance: 1000.0,
	}
	body, _ := json.Marshal(account)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createdAccount bank.Account
	err := json.Unmarshal(w.Body.Bytes(), &createdAccount)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.Owner, createdAccount.Owner)
	assert.Equal(suite.T(), account.InitialBalance, createdAccount.Balance)
}

func (suite *BankRestAPITestSuite) TestListAccountsHandler() {
	account := bank.Account{
		Owner:   "Alex Camara",
		Balance: 500.0,
	}
	_, err := suite.bankStore.CreateAccount(account.Owner, account.Balance)
	assert.NoError(suite.T(), err)

	req, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var accounts []bank.Account
	err = json.Unmarshal(w.Body.Bytes(), &accounts)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(accounts))
	assert.Equal(suite.T(), account.Owner, accounts[0].Owner)
}

func (suite *BankRestAPITestSuite) TestGetAccountByIDHandler() {
	account := bank.Account{
		Owner:   "Alex Camara",
		Balance: 3000.0,
	}
	createdAccount, err := suite.bankStore.CreateAccount(account.Owner, account.Balance)
	assert.NoError(suite.T(), err)

	req, _ := http.NewRequest(http.MethodGet, "/accounts/"+createdAccount.ID, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var retrievedAccount bank.Account
	err = json.Unmarshal(w.Body.Bytes(), &retrievedAccount)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdAccount.ID, retrievedAccount.ID)
}

func (suite *BankRestAPITestSuite) TestPerformTransactionHandler() {
	account := bank.Account{
		Owner:   "Alex Camara",
		Balance: 2000.0,
	}
	createdAccount, err := suite.bankStore.CreateAccount(account.Owner, account.Balance)
	assert.NoError(suite.T(), err)

	transaction := bank.Transaction{
		Type:   "deposit",
		Amount: 500.0,
	}
	transactionRequestBody, _ := json.Marshal(transaction)

	req, _ := http.NewRequest(http.MethodPost, "/accounts/"+createdAccount.ID+"/transactions", bytes.NewBuffer(transactionRequestBody))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createdTransaction bank.Transaction
	err = json.Unmarshal(w.Body.Bytes(), &createdTransaction)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), transaction.Amount, createdTransaction.Amount)
}

func (suite *BankRestAPITestSuite) TestTransferFundsHandler() {
	account1 := bank.Account{
		Owner:   "Alex Camara",
		Balance: 1500.0,
	}
	account2 := bank.Account{
		Owner:   "Donald Trump",
		Balance: 1000.0,
	}

	account1Created, err := suite.bankStore.CreateAccount(account1.Owner, account1.Balance)
	assert.NoError(suite.T(), err)
	account2Created, err := suite.bankStore.CreateAccount(account2.Owner, account2.Balance)
	assert.NoError(suite.T(), err)

	transfer := transferRequest{
		FromAccountID: account1Created.ID,
		ToAccountID:   account2Created.ID,
		Amount:        200.0,
	}
	transferRequestBody, _ := json.Marshal(transfer)

	req, _ := http.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(transferRequestBody))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	account1Updated, err := suite.bankStore.GetAccountByID(account1Created.ID)
	assert.NoError(suite.T(), err)
	account2Updated, err := suite.bankStore.GetAccountByID(account2Created.ID)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 1300.0, account1Updated.Balance)
	assert.Equal(suite.T(), 1200.0, account2Updated.Balance)
}

func (suite *BankRestAPITestSuite) TestGetTransactionsByAccountIDHandler() {
	account := bank.Account{
		Owner:   "Alex Camara",
		Balance: 1000.0,
	}
	createdAccount, err := suite.bankStore.CreateAccount(account.Owner, account.Balance)
	assert.NoError(suite.T(), err)

	transaction := bank.Transaction{
		Type:   "deposit",
		Amount: 300.0,
	}
	_, err = suite.bankStore.PerformTransaction(createdAccount.ID, transaction.Type, transaction.Amount)
	assert.NoError(suite.T(), err)

	req, _ := http.NewRequest(http.MethodGet, "/accounts/"+createdAccount.ID+"/transactions", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var transactions []bank.Transaction
	err = json.Unmarshal(w.Body.Bytes(), &transactions)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(transactions))
	assert.Equal(suite.T(), transaction.Amount, transactions[0].Amount)
}

func TestBankRestAPITestSuite(t *testing.T) {
	suite.Run(t, new(BankRestAPITestSuite))
}
