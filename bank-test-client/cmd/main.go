package main

import (
	"bank-test-client/tester"
	"fmt"
)

func main() {
	accountRepo := tester.NewAccountsRepo()
	for {
		tester.ShowMenu()
		choice := tester.GetMenuChoice()
		switch choice {
		case 1:
			tester.CreateAccount(accountRepo)
		case 2:
			tester.ListAccountsAndCompare(accountRepo)
		case 3:
			tester.RetrieveAccountDetails(accountRepo)
		case 4:
			tester.CreateTransaction(accountRepo)
		case 5:
			tester.GetAccountTransactions(accountRepo)
		case 6:
			tester.TransferFunds(accountRepo)
		case 7:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
