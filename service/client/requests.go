package client

type PurchaseRequest struct {
	Initiator    string `json:"initiator"`
	Amount       int    `json:"amount,omitempty"`
	Method       string `json:"method"`
	AccountId    int    `json:"account_id"`
	TargetNumber string `json:"target_number,omitempty"`
	DebitAccount string `json:"debit_account,omitempty"`
}

type UtilityPurchaseRequest struct {
	PurchaseRequest
	Provider      string
	AccountNumber string `json:"account_number"`
}

type VoucherPurchaseRequest struct {
	PurchaseRequest
	TargetAccountId int `json:"target_account_id,omitempty"`
}

type SubscriptionPurchaseRequest struct {
	PurchaseRequest
	SubscriptionTypeId int `json:"subscription_type_id,omitempty"`
}

type ProfileDetails struct {
	Name string
}

type SecurityQuestionRequest struct {
	QuestionId string `json:"question_id"`
	Answer     string
}

type EarningsWithdrawalRequest struct {
	Amount       int    `json:"amount"`
	Method       string `json:"method"`
	AccountId    int    `json:"account_id"`
	TargetNumber string `json:"target_number,omitempty"`
}

type NotificationRequest struct {
	Channel     string   `json:"channel"`
	Destination []string `json:"destination"`
	EventType   string   `json:"event_type"`
	Content     string   `json:"content"`
}

//TODO: Remove these once verified not needed
//func (r *AirtimePurchaseRequest) Marshal() ([]byte, error) {
//	request := map[string]interface{}{
//		"initiator":  r.Initiator,
//		"account_id": r.AccountId,
//		"amount":     r.Amount,
//		"method":     r.Method,
//	}
//
//	if r.TargetNumber != "" {
//		request["target_number"] = r.TargetNumber
//	}
//
//	if r.DebitAccount != "" {
//		request["debit_account"] = r.DebitAccount
//	}
//
//	return json.Marshal(r)
//}
//
//func (r *UtilityPurchaseRequest) Marshal() ([]byte, error) {
//	request := map[string]interface{}{
//		"initiator":      r.Initiator,
//		"account_id":     r.AccountId,
//		"amount":         r.Amount,
//		"method":         r.Method,
//		"provider":       r.Provider,
//		"account_number": r.AccountNumber,
//	}
//
//	return json.Marshal(request)
//}
