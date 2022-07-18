package client

import (
	"USSD.sidooh/cache"
	"USSD.sidooh/logger"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
	return api.baseUrl + endpoint
}

func (api *ApiClient) send(data interface{}) error {
	//TODO: Can we encode the data for security purposes and decode when necessary? Same to response logging...
	logger.ServiceLog.Println("API_REQ: ", api.request)
	start := time.Now()
	response, err := api.client.Do(api.request)
	if err != nil {
		logger.ServiceLog.Error("Error sending request to API endpoint: ", err)
		return err
	}
	// Close the connection to reuse it
	defer response.Body.Close()
	logger.ServiceLog.Println("API_RES - raw: ", response, time.Since(start))

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ServiceLog.Error("Couldn't parse response body: ", err)
	}
	logger.ServiceLog.Println("API_RES - body: ", string(body))

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
	if token := cache.Instance.Get("token"); token != nil /*&& !token.IsExpired()*/ {
		fmt.Println(token.Value(), token.IsExpired())
		// TODO: Check if token has expired since we should be able to decode it
		api.baseRequest(method, endpoint, body).request.Header.Add("Authorization", "Bearer "+token.Value())
	} else {
		api.EnsureAuthenticated()

		//TODO: What will happen to client if cache fails to store token? E.g. when account srv is not reachable?
		// TODO: Can we even just use a global Var?
		token = cache.Instance.Get("token")
		api.baseRequest(method, endpoint, body).request.Header.Add("Authorization", "Bearer "+token.Value())
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

	err := api.baseRequest(http.MethodPost, os.Getenv("ACCOUNTS_URL")+"/users/signin", bytes.NewBuffer(data)).send(response)
	if err != nil {
		return err
	}

	if cache.Instance != nil {
		cache.Instance.Set("token", response.Token, 14*time.Minute)
	}

	return nil
}
