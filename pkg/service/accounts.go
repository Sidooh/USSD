package service

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service/client"
	"fmt"
	"strconv"
)

func FetchAccount(phone string) (*client.Account, error) {
	account, err := accountsClient.GetAccountWithUser(phone)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch account", err)
		return nil, err
	}

	// TODO: make into goroutine
	// TODO: Making into goroutine means we have to set the details into session Vars somehow... so how???
	// 	Should we pass in vars map to func?
	/*go */
	func() {
		if account != nil && account.Id != 0 {
			// Check subscription
			subscription, err := FetchSubscription(strconv.Itoa(account.Id))
			if err != nil {
				logger.ServiceLog.Error("Failed to fetch user subscription: ", err)
			}

			account.Subscription = subscription

			// Check Voucher
			vouchers, err := paymentsClient.GetVouchers(strconv.Itoa(account.Id))
			if err != nil {
				logger.ServiceLog.Error("Failed to fetch voucher balances: ", err)
			}

			account.Vouchers = vouchers

			// Check Pin
			account.HasPin = CheckHasPin(strconv.Itoa(account.Id))
		}
	}()

	return account, nil
}

func FetchInvite(phone string) (*client.Invite, error) {
	// Check invite existence
	invite, err := accountsClient.CheckInvite(phone)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func CheckAccount(phone string) (*client.Account, error) {
	account, err := accountsClient.GetAccount(phone)
	if err != nil {
		logger.ServiceLog.Error("Failed to check account: ", err)
		return nil, err
	}

	return account, nil
}

// TODO: Add for Id only - may be faster when id is known?
func CheckAccountByIdOrPhone(search string) (*client.Account, error) {
	account, err := accountsClient.GetAccountByIdOrPhone(search)
	if err != nil {
		logger.ServiceLog.Error("Failed to search account: ", err)
		return nil, err
	}

	return account, nil
}

func InviteOrAccountExists(phone string) bool {
	// Check account existence
	account, err := accountsClient.GetAccount(phone)
	if err != nil {
		logger.ServiceLog.Error("Failed to check invite/account - account: ", err)
	}

	if account != nil {
		return true
	}

	// Check invite existence
	invite, err := accountsClient.CheckInvite(phone)
	if err != nil {
		logger.ServiceLog.Error("Failed to check invite/account - invite: ", err)
	}

	if invite != nil {
		return true
	}

	return false
}

func CheckPin(id string, pin string) bool {
	valid, err := accountsClient.CheckPin(id, pin)
	if err != nil {
		return false
	}

	return valid
}

func CheckHasPin(id string) bool {
	valid, err := accountsClient.CheckHasPin(id)
	if err != nil {
		return false
	}

	return valid
}

func CheckHasSecurityQuestions(id string) bool {
	valid, err := accountsClient.CheckHasSecurityQuestions(id)
	if err != nil {
		return false
	}

	return valid
}

func CreateAccount(phone string, inviteCode interface{}) (*client.Account, error) {
	account, err := accountsClient.CreateAccount(phone, inviteCode)
	if err != nil {
		return nil, err
	}

	go func() {
		if account.InviterId != 0 {
			inviter, err := CheckAccountByIdOrPhone(strconv.Itoa(account.InviterId))
			if err != nil {
				logger.ServiceLog.Error("Failed to fetch invite account: ", err)
				return
			}

			message := fmt.Sprintf("Congratulations! %s has "+
				"successfully accessed Sidooh using your invite code. "+
				"Show them how to buy airtime from Sidooh so as to unlock your earnings."+
				"The more friends you invite to Sidooh, the more you earn.", account.Phone)
			request := client.NotificationRequest{
				Channel:     "SMS",
				Destination: []string{inviter.Phone},
				EventType:   "REFERRAL_JOINED", //TODO: Change notify referral types to invite
				Content:     message,
			}

			Notify(&request)
		}
	}()

	return account, nil
}

func CreateInvite(id string, phone string) (*client.Invite, error) {
	invite, err := accountsClient.CreateInvite(id, phone)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func SetPin(id string, pin string) bool {
	valid, err := accountsClient.SetPin(id, pin)
	if err != nil {
		return false
	}

	return valid
}

func UpdateProfile(id string, details client.ProfileDetails) (*client.User, error) {
	user, err := accountsClient.UpdateProfile(id, details)
	if err != nil {
		return user, err
	}

	return user, nil
}

var securityQuestions []client.SecurityQuestion

func FetchSecurityQuestions() ([]client.SecurityQuestion, error) {
	if len(securityQuestions) > 0 {
		return securityQuestions, nil
	}

	questions, err := accountsClient.FetchSecurityQuestions()
	if err != nil {
		return nil, err
	}

	securityQuestions = questions

	return questions, nil
}

func SetSecurityQuestions(id string, answers map[string]string) error {
	var results []interface{}

	for i, answer := range answers {
		res, err := accountsClient.SetSecurityQuestion(id, client.SecurityQuestionRequest{
			QuestionId: i,
			Answer:     answer,
		})
		if err != nil {
			return err
		} else {
			results = append(results, res)
		}
	}

	logger.ServiceLog.Println("Security Questions Request: ", results)
	return nil
}

func FetchUserSecurityQuestions(id string) ([]client.UserSecurityQuestion, error) {
	questions, err := accountsClient.FetchUserSecurityQuestions(id)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func CheckSecurityQuestionAnswers(id string, answers map[string]string) bool {
	var valid = false

	for i, answer := range answers {
		var res map[string]bool
		res, err := accountsClient.CheckSecurityQuestionAnswers(id, client.SecurityQuestionRequest{
			QuestionId: i,
			Answer:     answer,
		})
		if err != nil {
			return false
		} else {
			valid = res["message"]
		}
	}

	return valid
}

func FetchEarningBalances(id string) ([]client.EarningAccount, error) {
	earnings, err := productsClient.FetchAccountEarnings(id)
	if err != nil {
		return nil, err
	}

	return earnings, nil
}

func FetchSavingBalances(id string) ([]client.SavingAccount, error) {
	earnings, err := savingsClient.FetchAccountSavings(id)
	if err != nil {
		return nil, err
	}

	return earnings, nil
}

func RequestEarningsWithdrawal(request *client.EarningsWithdrawalRequest) error {
	err := productsClient.WithdrawEarnings(request)
	if err != nil {
		logger.ServiceLog.Error("Failed to withdraw earnings: ", err)
		return err
	}

	return nil
}
