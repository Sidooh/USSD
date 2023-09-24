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

func BuyFloat(id string, request client.FloatPurchaseRequest) {
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
