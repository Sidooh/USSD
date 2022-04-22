package client

import (
	"USSD.sidooh/logger"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type ApiClient struct {
	client  *http.Client
	request *http.Request
	baseUrl string
}

type AuthResponse struct {
	Token string `json:"token"`
}

var (
	cache = ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](14 * time.Minute),
	)
)

func (api *ApiClient) init(baseUrl string) {
	api.client = &http.Client{Timeout: 10 * time.Second}
	api.baseUrl = baseUrl
}

func (api *ApiClient) getUrl(endpoint string) string {
	return api.baseUrl + endpoint
}

func (api *ApiClient) send(data interface{}) error {
	response, err := api.client.Do(api.request)
	if err != nil {
		logger.ServiceLog.Fatalf("Error sending request to API endpoint. %+v", err)
	}
	// Close the connection to reuse it
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ServiceLog.Fatalf("Couldn't parse response body. %+v", err)
	}

	logger.ServiceLog.Println(response)
	if response.StatusCode != 200 && response.StatusCode != 401 && response.StatusCode != 404 {
		if response.StatusCode < 500 {
			var errorMessage map[string][]map[string]string
			err = json.Unmarshal(body, &errorMessage)

			return errors.New(errorMessage["errors"][0]["message"])
		}

		return errors.New(string(body))
	}

	if response.StatusCode == 404 {
		return errors.New(string(body))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.ServiceLog.Fatalf("Failed to unmarshall body. %+v", err)
	}

	return nil
}

func (api *ApiClient) setDefaultHeaders() {
	api.request.Header.Add("Accept", "application/json")
	api.request.Header.Add("Content-Type", `application/json`)
}

func (api *ApiClient) baseRequest(method string, endpoint string, body io.Reader) *ApiClient {
	endpoint = api.getUrl(endpoint)
	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		logger.ServiceLog.Fatalf("error creating HTTP request: %v", err)
	}

	api.request = request
	api.setDefaultHeaders()

	logger.ServiceLog.Println(body)

	return api
}

func (api *ApiClient) newRequest(method string, endpoint string, body io.Reader) *ApiClient {
	if token := cache.Get("token"); token != nil {
		api.baseRequest(method, endpoint, body).request.Header.Add("Authorization", "Bearer "+token.Value())
	} else {
		api.EnsureAuthenticated()

		token = cache.Get("token")
		api.baseRequest(method, endpoint, body).request.Header.Add("Authorization", "Bearer "+token.Value())
	}

	return api
}

func (api *ApiClient) EnsureAuthenticated() {
	values := map[string]string{"email": "aa@a.a", "password": "12345678"}
	jsonData, err := json.Marshal(values)

	err = api.Authenticate(jsonData)
	if err != nil {
		logger.ServiceLog.Fatalf("error authenticating: %v", err)
	}
}

func (api *ApiClient) Authenticate(data []byte) error {
	var response = new(AuthResponse)

	err := api.baseRequest(http.MethodPost, "/users/signin", bytes.NewBuffer(data)).send(response)
	if err != nil {
		return err
	}

	cache.Set("token", response.Token, 14*time.Minute)
	go func() {
		for {
			time.Sleep(14 * time.Minute)
			cache.Delete("token")
		}
	}()

	return nil
}
