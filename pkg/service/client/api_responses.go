package client

type Account struct {
	Id           int    `json:"id,omitempty"`
	Phone        string `json:"phone"`
	Active       bool   `json:"active"`
	InviterId    int    `json:"inviter_id"`
	User         `json:"user"`
	Vouchers     []Voucher
	Subscription Subscription
	HasPin       bool
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Voucher struct {
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
	VoucherType struct {
		Name string `json:"name"`
	}
}

type Invite struct {
	Id      int      `json:"id"`
	Phone   string   `json:"phone"`
	Status  string   `json:"status"`
	Inviter *Account `json:"inviter"`
}

type UtilityAccount struct {
	Id            int
	Provider      string
	AccountNumber string `json:"account_number"`
}

type SubscriptionType struct {
	Id       int
	Title    string
	Price    int
	Duration int
	Active   bool
}

type Subscription struct {
	Id        int
	Status    string `json:"status"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type EarningRate struct {
	Type  string
	Value float64
}
