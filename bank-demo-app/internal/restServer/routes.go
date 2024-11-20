package restServer

import (
	"net/http"
)

func InitRestRoutes(bankStore BankStore) Routes {
	serverRoutes := Routes{
		// Route to get if server is up.
		{
			Method:  http.MethodGet,
			Pattern: "/status",
			Handler: statusHandler,
		},
		// Create a new bank account with an initial balance.
		{
			Method:  http.MethodPost,
			Pattern: "/accounts",
			Handler: createAccountHandler(bankStore),
		},
		// Retrieve a list of all bank accounts.
		{
			Method:  http.MethodGet,
			Pattern: "/accounts",
			Handler: listAccountsHandler(bankStore),
		},
		// Retrieve details of a specific account by ID.
		{
			Method:  http.MethodGet,
			Pattern: "/accounts/:id",
			Handler: getAccountByIDHandler(bankStore),
		},
		// Create a deposit or withdrawal transaction for a specific account.
		{
			Method:  http.MethodPost,
			Pattern: "/accounts/:id/transactions",
			Handler: performTransactionHandler(bankStore),
		},
		// Retrieve all transactions sassociated with a specific account.
		{
			Method:  http.MethodGet,
			Pattern: "/accounts/:id/transactions",
			Handler: getTransactionsByAccountIDHandler(bankStore),
		},
		// Transfer funds from one account to another.
		{
			Method:  http.MethodPost,
			Pattern: "/transfer",
			Handler: transferFundsHandler(bankStore),
		},
	}
	return serverRoutes
}
