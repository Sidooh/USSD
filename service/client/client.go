package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type ApiClient struct {
	client  *http.Client
	request *http.Request
}

type AuthResponse struct {
	Token string `json:"token"`
}

var (
	cache = ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](14 * time.Minute),
	)
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

func (api *ApiClient) setDefaultHeaders() {
	api.request.Header.Add("Accept", "application/json")
	api.request.Header.Add("Content-Type", `application/json`)
}

func (api *ApiClient) baseRequest(method string, endpoint string, body io.Reader) *ApiClient {
	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	api.request = request
	api.setDefaultHeaders()

	return api
}

func (api *ApiClient) newRequest(method string, endpoint string, body io.Reader) *ApiClient {
	if token := cache.Get("token"); token != nil {
		fmt.Println("Exp", token.ExpiresAt(), token.TTL())
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
		log.Fatalf("error authenticating: %v", err)
	}
}

func (api *ApiClient) Authenticate(data []byte) error {
	var response = new(AuthResponse)

	err := api.baseRequest(http.MethodPost, getUrl("/users/signin"), bytes.NewBuffer(data)).send(response)
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
