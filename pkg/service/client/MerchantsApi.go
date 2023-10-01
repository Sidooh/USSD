package client

import (
	"USSD.sidooh/pkg/cache"
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type MerchantsApiClient struct {
	ApiClient
}

type CountyApiResponse struct {
	ApiResponse
	Data *[]County `json:"data"`
}

type SubCountyApiResponse struct {
	ApiResponse
	Data *[]SubCounty `json:"data"`
}

type WardApiResponse struct {
	ApiResponse
	Data *[]Ward `json:"data"`
}

type LandmarkApiResponse struct {
	ApiResponse
	Data *[]Landmark `json:"data"`
}

type MerchantApiResponse struct {
	ApiResponse
	Data *Merchant `json:"data"`
}

type MpesaStoreAccountApiResponse struct {
	ApiResponse
	Data *[]MpesaStoreAccount `json:"data"`
}

type MerchantEarningAccountApiResponse struct {
	ApiResponse
	Data *[]MerchantEarningAccount `json:"data"`
}

func InitMerchantClient() *MerchantsApiClient {
	client := MerchantsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_MERCHANTS_API_URL"))
	return &client
}

func (m *MerchantsApiClient) GetMerchantByAccount(accountId string) (*Merchant, error) {
	var res = new(MerchantApiResponse)

	err := m.newRequest(http.MethodGet, "/merchants/account/"+accountId, nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (m *MerchantsApiClient) GetCounties() (*[]County, error) {
	var res = new(CountyApiResponse)

	counties, err := cache.Get[[]County]("counties")
	if err == nil && counties != nil {
		return counties, nil
	}

	err = m.newRequest(http.MethodGet, "/counties", nil).send(res)
	if err != nil {
		return nil, err
	}

	cache.Set("counties", res.Data, 360*24*time.Hour)

	return res.Data, nil
}

func (m *MerchantsApiClient) GetSubCounties(county string) (*[]SubCounty, error) {
	var res = new(SubCountyApiResponse)

	subCounties, err := cache.Get[[]SubCounty]("sub_counties-" + county)
	if err == nil && subCounties != nil {
		return subCounties, nil
	}

	err = m.newRequest(http.MethodGet, "/counties/"+county, nil).send(res)
	if err != nil {
		return nil, err
	}

	cache.Set("sub_counties-"+county, res.Data, 360*24*time.Hour)

	return res.Data, nil
}

func (m *MerchantsApiClient) GetWards(county string, subCounty string) (*[]Ward, error) {
	var res = new(WardApiResponse)

	wards, err := cache.Get[[]Ward]("wards-" + subCounty)
	if err == nil && wards != nil {
		return wards, nil
	}

	err = m.newRequest(http.MethodGet, "/counties/"+county+"/sub-counties/"+subCounty, nil).send(res)
	if err != nil {
		return nil, err
	}

	cache.Set("wards-"+subCounty, res.Data, 360*24*time.Hour)

	return res.Data, nil
}

func (m *MerchantsApiClient) GetLandmarks(county string, subCounty string, ward string) (*[]Landmark, error) {
	var res = new(LandmarkApiResponse)

	landmarks, err := cache.Get[[]Landmark]("landmarks-" + ward)
	if err == nil && landmarks != nil {
		return landmarks, nil
	}

	err = m.newRequest(http.MethodGet, "/counties/"+county+"/sub-counties/"+subCounty+"/wards/"+ward, nil).send(res)
	if err != nil {
		return nil, err
	}

	cache.Set("landmarks-"+subCounty, res.Data, 360*24*time.Hour)

	return res.Data, nil
}

func (m *MerchantsApiClient) CreateMerchant(request MerchantKYCDetails) (*Merchant, error) {
	var res = new(MerchantApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = m.newRequest(http.MethodPost, "/merchants", dataBytes).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (m *MerchantsApiClient) UpdateKYBData(id string, request MerchantKYBDetails) (*Merchant, error) {
	var res = new(MerchantApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = m.newRequest(http.MethodPost, "/merchants/"+id+"/kyb", dataBytes).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (m *MerchantsApiClient) BuyFloat(id string, request FloatPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = m.newRequest(http.MethodPost, "/merchants/"+id+"/buy-float", dataBytes).send(nil)
	return err
}

func (m *MerchantsApiClient) GetMerchantByIdNumber(id string) (*Merchant, error) {
	var res = new(MerchantApiResponse)

	err := m.newRequest(http.MethodGet, "/merchants/id-number/"+id, nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (p *MerchantsApiClient) GetMpesaStoreAccounts(id string) (*[]MpesaStoreAccount, error) {
	res := new(MpesaStoreAccountApiResponse)

	err := p.newRequest(http.MethodGet, "/merchants/"+id+"/mpesa-store-accounts", nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (p *MerchantsApiClient) GetEarningAccounts(id string) (*[]MerchantEarningAccount, error) {
	res := new(MerchantEarningAccountApiResponse)

	err := p.newRequest(http.MethodGet, "/earning-accounts/merchant/"+id, nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}
