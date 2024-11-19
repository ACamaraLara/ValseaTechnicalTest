package tester

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Sends create account request and stores it's response inside local repo to let the user perform checks and transactions without knowing the .
func CreateAccount(accountRepo AccountsRepo) {
	fmt.Println("\nCreate a New Bank Account:")
	fmt.Print("Enter account owner's name: ")
	owner := InputString()

	fmt.Print("Enter initial balance: ")
	initialBalance := InputFloat()

	account := map[string]interface{}{
		"owner":           owner,
		"initial_balance": initialBalance,
	}

	resp, err := makeRequest(http.MethodPost, "/accounts", account)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Error: Expected StatusCreated code but received %d\n", resp.StatusCode)
		printResponseBody(resp.Body)
		return
	}

	var createdAccount Account
	err = json.NewDecoder(resp.Body).Decode(&createdAccount)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return
	}

	fmt.Println("Created account:")
	printResponseJson(createdAccount)

	// Adding the new account to the accountRepo, keyed by the owner's name to use it in future tansactions.
	accountRepo[createdAccount.Owner] = &createdAccount
}

// Gets all the accounts stored inside the REST server and checks if are the same than the ones that
// have been created by this testing application.
func ListAccountsAndCompare(accountRepo AccountsRepo) {
	fmt.Println("\nList all accounts received from" +
		" server and compare them with the stored ones:")

	resp, err := makeRequest(http.MethodGet, "/accounts", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Expected StatusOK code but received %d\n", resp.StatusCode)
		printResponseBody(resp.Body)
		return
	}

	var receivedAccounts []Account
	err = json.NewDecoder(resp.Body).Decode(&receivedAccounts)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return
	}

	fmt.Println("Received accounts from server:")
	printResponseJson(receivedAccounts)

	// Check if the accounts from the API match those in the local accountRepo
	matches := true
	for _, apiAccount := range receivedAccounts {
		storedAccount, exists := accountRepo[apiAccount.Owner]
		if !exists {
			fmt.Printf("Account with owner %s not found in repo.\n", apiAccount.Owner)
			matches = false
			continue
		}
		// Compare the account details
		if *storedAccount != apiAccount {
			fmt.Printf("Account details for %s do not match. API: %+v, Repo: %+v\n", apiAccount.Owner, apiAccount, storedAccount)
			matches = false
		}
	}

	if matches {
		fmt.Println("All accounts match between the API and the repo.")
	} else {
		fmt.Println("There are discrepancies between the API and the repo accounts.")
	}

}

// Sends create account requests and stores it's response inside local repo to let the user perform checks and transactions without knowing the .
func RetrieveAccountDetails(accountRepo AccountsRepo) {
	fmt.Println("\nCreate a New Bank Account:")
	fmt.Println("Enter account owner's name: ")
	owner := InputString()
	storedAccount, exists := accountRepo[owner]
	if !exists {
		fmt.Printf("Account with owner %s not found in repo.\n", owner)
		return
	}

	resp, err := makeRequest(http.MethodGet, "/accounts/"+storedAccount.ID, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Expected StatusOK code but received %d\n", resp.StatusCode)
		printResponseBody(resp.Body)
		return
	}

	var receivedAccount Account
	err = json.NewDecoder(resp.Body).Decode(&receivedAccount)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return
	}

	fmt.Println("Requested account details:")
	printResponseJson(receivedAccount)
}

// CreateTransaction handles creating a transaction for a specific account.
func CreateTransaction(accountRepo AccountsRepo) {
	fmt.Println("\nCreate a Transaction for an Existing Account:")

	storedAccount, owner := getAccountByOwner("Enter account owner's name: ", accountRepo)
	if storedAccount == nil {
		fmt.Printf("Account with owner %s not found in repo.\n", owner)
		return
	}

	fmt.Println("Choose transaction type (deposit/withdrawal): ")
	transactionType := InputString()

	fmt.Print("Enter transaction amount: ")
	amount := InputFloat()

	transaction := map[string]interface{}{
		"type":   transactionType,
		"amount": amount,
	}

	resp, err := makeRequest(http.MethodPost, "/accounts/"+storedAccount.ID+"/transactions", transaction)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Error: Expected StatusCreated code but received %d\n", resp.StatusCode)
		printResponseBody(resp.Body)
		return
	}

	updateAccountAfterTransaction(storedAccount, transactionType, amount)

	var newTransaction Transaction
	err = json.NewDecoder(resp.Body).Decode(&newTransaction)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return
	}
	fmt.Println("Transaction successfull, details:")
	printResponseJson(newTransaction)
}

func updateAccountAfterTransaction(account *Account, trxType string, amount float64) {
	if trxType == "deposit" {
		account.Deposit(amount)
	} else if trxType == "withdrawal" {
		account.Withdraw(amount)
	}
}

// GetAccountTransactions retrieves all transactions for a specific account.
func GetAccountTransactions(accountRepo AccountsRepo) {
	fmt.Println("\nRetrieve Transactions for a Specific Account:")

	storedAccount, owner := getAccountByOwner("Enter account owner's name: ", accountRepo)
	if storedAccount == nil {
		fmt.Printf("Account with owner %s not found in repo.\n", owner)
		return
	}

	resp, err := makeRequest(http.MethodGet, "/accounts/"+storedAccount.ID+"/transactions", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Expected StatusOK code but received %d\n", resp.StatusCode)
		printResponseBody(resp.Body)
		return
	}

	var transactions []Transaction
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return
	}

	if len(transactions) == 0 {
		fmt.Println("No transactions found for this account.")
		return
	}

	fmt.Printf("\nTransactions for account %s: \n", storedAccount.Owner)
	printResponseJson(transactions)
}

// TransferFunds allows transferring funds between two accounts.
func TransferFunds(accountRepo AccountsRepo) {
	fmt.Println("Transfer Funds Between Accounts:")

	fromAccount, fromOwner := getAccountByOwner("Enter sender account owner's name: ", accountRepo)
	if fromAccount == nil {
		fmt.Printf("Account with owner %s not found in repo.\n", fromOwner)
		return
	}

	toAccount, toOwner := getAccountByOwner("Enter receiver account owner's name: ", accountRepo)
	if toAccount == nil {
		fmt.Printf("Account with owner %s not found in repo.\n", toOwner)
		return
	}

	fmt.Print("Enter transfer amount: ")
	amount := InputFloat()

	transferRequest := map[string]interface{}{
		"from_account_id": fromAccount.ID,
		"to_account_id":   toAccount.ID,
		"amount":          amount,
	}

	resp, err := makeRequest(http.MethodPost, "/transfer", transferRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Expected StatusOK code but received %d\n", resp.StatusCode)
		printResponseBody(resp.Body)
		return
	}

	updateAccountAfterTransaction(fromAccount, "withdrawal", amount)
	updateAccountAfterTransaction(toAccount, "deposit", amount)

	printResponseBody(resp.Body)
	fmt.Println("Transfer successful!")
}
