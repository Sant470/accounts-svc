package apptypes

type AccountsRequest struct {
	AccountID      string  `json:"account_id"`
	InitialBalance float64 `json:"initial_balance"`
}

type TransactionRequest struct {
	SourceAccountID      string  `json:"source_account_id"` // add tags for validation as well
	DestinationAccountID string  `json:"destination_account_id"`
	Amount               float64 `json:"amount"`
}
