package restServer

// request struct for creating an account
type createAccountRequest struct {
	Owner          string  `json:"owner"`
	InitialBalance float64 `json:"initial_balance"`
}

// request struct for creating a transaction
type createTransactionRequest struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

// request struct for transferring funds
type transferRequest struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}
