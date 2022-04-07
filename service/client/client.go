package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ApiClient struct {
	client  *http.Client
	request *http.Request
}

type AuthResponse struct {
	Token string `json:"token"`
}

var (
	token  = ""
	ticker = time.NewTicker(14 * time.Minute)
	quit   = make(chan struct{})
)

func (api *ApiClient) init() {
	api.client = &http.Client{Timeout: 10 * time.Second}
}

func (api *ApiClient) send(data interface{}) error {
	response, err := api.client.Do(api.request)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}
	// Close the connection to reuse it
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	return nil
}

func (api *ApiClient) newRequest(method string, endpoint string, body io.Reader) *ApiClient {
	api.init()

	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", `application/json`)

	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}

	api.request = request
	return api
}

func (api ApiClient) EnsureAuthenticated() {
	if token == "" {
		values := map[string]string{"email": "aa@a.a", "password": "12345678"}
		jsonData, err := json.Marshal(values)

		err = api.Authenticate(jsonData)
		if err != nil {
			log.Fatalf("error authenticating: %v", err)
		}
	}
}

func (a *ApiClient) Authenticate(data []byte) error {
	var response = new(AuthResponse)

	err := a.newRequest(http.MethodPost, getUrl("/users/signin"), bytes.NewBuffer(data)).send(response)
	if err != nil {
		return err
	}

	token = response.Token

	//TODO: Test this worker for unsetting token every so often
	//
	//go func() {
	//	for {
	//		select {
	//		case <-ticker.C:
	//			token = ""
	//		case <-quit:
	//			ticker.Stop()
	//			return
	//		}
	//	}
	//}()

	return nil
}
