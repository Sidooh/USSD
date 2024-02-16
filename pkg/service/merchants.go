package service

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service/client"
)

func FetchCounties() (counties *[]client.County, err error) {
	counties, err = merchantsClient.GetCounties()
	return
}

func FetchSubCounties(county string) (subCounties *[]client.SubCounty, err error) {
	subCounties, err = merchantsClient.GetSubCounties(county)
	return
}

func FetchWards(county, subCounty string) (wards *[]client.Ward, err error) {
	wards, err = merchantsClient.GetWards(county, subCounty)
	return
}

func FetchLandmarks(county, subCounty, ward string) (landmarks *[]client.Landmark, err error) {
	landmarks, err = merchantsClient.GetLandmarks(county, subCounty, ward)
	return
}

func CreateMerchant(request client.MerchantKYCDetails) (merchant *client.Merchant, err error) {
	merchant, err = merchantsClient.CreateMerchant(request)
	return
}

func UpdateMerchant(id string, request client.MerchantKYBDetails) (merchant *client.Merchant, err error) {
	merchant, err = merchantsClient.UpdateKYBData(id, request)
	return
}

func GetMerchantByAccount(accountId string) (merchant *client.Merchant, err error) {
	merchant, err = merchantsClient.GetMerchantByAccount(accountId)
	return
}

func BuyFloat(id string, request client.MerchantMpesaFloatPurchaseRequest) {
	err := merchantsClient.BuyFloat(id, request)
	if err != nil {
		logger.ServiceLog.Error("Failed to buy float: ", err)
	}
}

func MerchantIdNumberExists(id string) bool {
	merchant, _ := merchantsClient.GetMerchantByIdNumber(id)
	return merchant != nil
}

func FetchMpesaStoreAccounts(merchantId string) ([]client.MpesaStoreAccount, error) {
	accounts, err := merchantsClient.GetMpesaStoreAccounts(merchantId)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch mpesa store accounts: ", err)
		return nil, err
	}

	return *accounts, nil
}

func FetchEarningAccounts(merchantId string) ([]client.MerchantEarningAccount, error) {
	accounts, err := merchantsClient.GetEarningAccounts(merchantId)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch earning accounts: ", err)
		return nil, err
	}

	return accounts, nil
}

func FetchTransactions(accountIds, days string) ([]client.Transaction, error) {
	transactions, err := merchantsClient.FetchTransactions(accountIds, days)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch transactions: ", err)
		return nil, err
	}

	return *transactions, nil
}

func WithdrawEarnings(id string, request client.MerchantWithdrawalRequest) error {
	err := merchantsClient.WithdrawEarnings(id, request)
	if err != nil {
		logger.ServiceLog.Error("Failed to withdraw earnings: ", err)
	}

	return err
}

func MpesaWithdrawal(id string, request client.MerchantMpesaWithdrawalRequest) error {
	err := merchantsClient.MpesaWithdrawal(id, request)
	if err != nil {
		logger.ServiceLog.Error("Failed to withdraw mpesa: ", err)
	}

	return err
}

func VoucherPurchase(id string, request client.MerchantMpesaWithdrawalRequest) error {
	err := merchantsClient.VoucherPurchase(id, request)
	if err != nil {
		logger.ServiceLog.Error("Failed to purchase voucher: ", err)
	}

	return err
}

func VoucherTransfer(id string, request client.MerchantFloatTransferRequest) error {
	err := merchantsClient.VoucherTransfer(id, request)
	if err != nil {
		logger.ServiceLog.Error("Failed to transfer voucher: ", err)
	}

	return err
}

func VoucherWithdraw(id string, request client.MerchantWithdrawalRequest) error {
	err := merchantsClient.VoucherWithdraw(id, request)
	if err != nil {
		logger.ServiceLog.Error("Failed to withdraw voucher: ", err)
	}

	return err
}
