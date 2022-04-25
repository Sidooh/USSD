package client

import (
	"encoding/json"
)

type AirtimePurchaseRequest struct {
	Initiator    string `json:"initiator"`
	Amount       int    `json:"amount"`
	Method       string `json:"method"`
	AccountId    int    `json:"account_id"`
	TargetNumber string `json:"target_number,omitempty"`
	DebitAccount string `json:"debit_account,omitempty"`
}

type UtilityPurchaseRequest struct {
	AirtimePurchaseRequest
	Provider      string
	AccountNumber string `json:"account_number"`
}

func (r *AirtimePurchaseRequest) Marshal() ([]byte, error) {
	request := map[string]interface{}{
		"initiator":  r.Initiator,
		"account_id": r.AccountId,
		"amount":     r.Amount,
		"method":     r.Method,
	}

	if r.TargetNumber != "" {
		request["target_number"] = r.TargetNumber
	}

	if r.DebitAccount != "" {
		request["debit_account"] = r.DebitAccount
	}

	return json.Marshal(r)
}

func (r *UtilityPurchaseRequest) Marshal() ([]byte, error) {
	request := map[string]interface{}{
		"initiator":      r.Initiator,
		"account_id":     r.AccountId,
		"amount":         r.Amount,
		"method":         r.Method,
		"provider":       r.Provider,
		"account_number": r.AccountNumber,
	}

	return json.Marshal(request)
}
