package client

import (
	"USSD.sidooh/pkg/cache"
	"USSD.sidooh/pkg/logger"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
	"time"
)

type ApiClient struct {
	client  *http.Client
	request *http.Request
	baseUrl string
}

type AuthResponse struct {
	Token string `json:"access_token"`
}

type ApiResponse struct {
	Result  int           `json:"result"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
	Errors  []interface{} `json:"errors"`
}

func (api *ApiClient) init(baseUrl string) {
	logger.ServiceLog.Println("Init client: ", baseUrl)

	// TODO: Review switching to http2
	api.client = &http.Client{Timeout: 10 * time.Second}
	api.baseUrl = baseUrl
}

func (api *ApiClient) getUrl(endpoint string) string {
	if strings.HasPrefix(endpoint, "http") {
		return endpoint
	}
	if !strings.HasPrefix(api.baseUrl, "http") {
		api.baseUrl = "https://" + api.baseUrl
	}
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	return api.baseUrl + endpoint
}

func (api *ApiClient) send(data interface{}) error {
	//TODO: Can we encode the data for security purposes and decode when necessary? Same to response logging...
	logger.ServiceLog.WithField("req", fmt.Sprint(api.request)).Println("API_REQ: ")
	start := time.Now()
	response, err := api.client.Do(api.request)
	if err != nil {
		logger.ServiceLog.Error("Error sending request to API endpoint: ", err)
		return err
	}

	// Close the connection to reuse it
	defer response.Body.Close()
	logger.ServiceLog.WithField("res", fmt.Sprint(response)).WithField("latency (ms)", time.Since(start).Milliseconds()).Println("API_RES - raw: ")

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.ServiceLog.Error("Couldn't parse response body: ", err)
	}

	logger.ServiceLog.WithField("body", string(body)).Println("API_RES - body: ")

	//TODO: Perform error handling in a better way
	if response.StatusCode != 200 && response.StatusCode != 201 && response.StatusCode != 401 &&
		response.StatusCode != 404 && response.StatusCode != 422 {
		if response.StatusCode < 500 {
			var errorMessage map[string][]map[string]string
			err = json.Unmarshal(body, &errorMessage)

			if len(errorMessage["errors"]) == 0 {
				var errorMessage map[string]string
				err = json.Unmarshal(body, &errorMessage)
				logger.ServiceLog.Error("API_ERR - body: ", errorMessage)

				return errors.New(errorMessage["message"])
			}

			return errors.New(errorMessage["errors"][0]["message"])
		}

		return errors.New(string(body))
	}

	if response.StatusCode == 404 {
		return errors.New(string(body))
	}

	//TODO: Deal with 401
	if response.StatusCode == 401 {
		logger.ServiceLog.Panic("Failed to authenticate.")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.ServiceLog.Error("Failed to unmarshal body: ", err)
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

	return api
}

func (api *ApiClient) newRequest(method string, endpoint string, body io.Reader) *ApiClient {
	if token, err := cache.GetString("token"); err == nil {
		// TODO: Check if token has expired since we should be able to decode it
		api.baseRequest(method, endpoint, body).request.Header.Add("Authorization", "Bearer "+token)
	} else {
		api.EnsureAuthenticated()

		//TODO: What will happen to client if cache fails to store token? E.g. when account srv is not reachable?
		// TODO: Can we even just use a global Var?
		token, _ = cache.GetString("token")
		api.baseRequest(method, endpoint, body).request.Header.Add("Authorization", "Bearer "+token)
	}

	return api
}

func (api *ApiClient) EnsureAuthenticated() {
	values := map[string]string{"email": "aa@a.a", "password": "12345678"}
	jsonData, err := json.Marshal(values)

	err = api.Authenticate(jsonData)
	if err != nil {
		logger.ServiceLog.Errorf("error authenticating: %v", err)
	}
}

func (api *ApiClient) Authenticate(data []byte) error {
	var response = new(AuthResponse)

	err := api.baseRequest(http.MethodPost, viper.GetString("SIDOOH_ACCOUNTS_API_URL")+"/users/signin", bytes.NewBuffer(data)).send(response)
	if err != nil {
		return err
	}

	cache.SetString("token", response.Token, 14*time.Minute)

	return nil
}
