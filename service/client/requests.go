package client

import (
	"encoding/json"
)

type AirtimePurchaseRequest struct {
	Initiator    string `json:"initiator"`
	Amount       int    `json:"amount"`
	Method       string `json:"method"`
	AccountId    int    `json:"account_id"`
	TargetNumber int    `json:"target_number"`
	DebitAccount int    `json:"debit_account"`
}

func (r *AirtimePurchaseRequest) Marshal() ([]byte, error) {
	request := map[string]interface{}{
		"initiator":  r.Initiator,
		"account_id": r.AccountId,
		"amount":     r.Amount,
		"method":     r.Method,
	}

	if r.TargetNumber != 0 {
		request["target_number"] = r.TargetNumber
	}

	if r.DebitAccount != 0 {
		request["debit_account"] = r.DebitAccount
	}

	return json.Marshal(request)
}
