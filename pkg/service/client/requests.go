package client

type PurchaseRequest struct {
	Initiator    string `json:"initiator"`
	Amount       int    `json:"amount,omitempty"`
	Method       string `json:"method,omitempty"`
	AccountId    int    `json:"account_id"`
	TargetNumber string `json:"target_number,omitempty"`
	DebitAccount string `json:"debit_account,omitempty"`
}

type UtilityPurchaseRequest struct {
	PurchaseRequest
	Provider      string `json:"provider"`
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

type MerchantPurchaseRequest struct {
	PurchaseRequest
	MerchantType   string `json:"merchant_type"`
	BusinessNumber string `json:"business_number"`
	AccountNumber  string `json:"account_number,omitempty"`
}

type ProfileDetails struct {
	Name string
}

type SecurityQuestionRequest struct {
	QuestionId string `json:"question_id"`
	Answer     string
}

type EarningsWithdrawalRequest struct {
	PurchaseRequest
}

type NotificationRequest struct {
	Channel     string   `json:"channel"`
	Destination []string `json:"destination"`
	EventType   string   `json:"event_type"`
	Content     string   `json:"content"`
}

type MerchantKYCDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IdNumber  string `json:"id_number"`
	AccountId int    `json:"account_id"`
}

type MerchantKYBDetails struct {
	BusinessName string `json:"business_name"`
	Landmark     string `json:"landmark"`
}
