package client

import (
	"USSD.sidooh/cache"
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
)

var client ApiClient

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGetUrl(t *testing.T) {
	// 1. https 2. http 3. slashes 4.

	testData := map[string]string{
		"url1":             "https://test.url/url1",
		"/url1":            "https://test.url/url1",
		"https://test.url": "https://test.url",
		"http://test.url":  "http://test.url",
	}

	// Tests bare base url
	client.baseUrl = "test.url"
	performUrlTest(t, testData)

	// Tests https base url
	client.baseUrl = "https://test.url"
	performUrlTest(t, testData)

	testData = map[string]string{
		"url1":             "http://test.url/url1",
		"/url1":            "http://test.url/url1",
		"https://test.url": "https://test.url",
		"http://test.url":  "http://test.url",
	}
	// Tests http base url
	client.baseUrl = "http://test.url"
	performUrlTest(t, testData)
}

func performUrlTest(t *testing.T, testData map[string]string) {
	urlReg := "[(http(s)?):\\/\\/(www\\.)?a-zA-Z0-9@:%._\\+~#=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)"

	for endpoint, expected := range testData {
		url := client.getUrl(endpoint)
		_, err := regexp.MatchString(urlReg, url)
		if err != nil {
			t.Errorf("getUrl = %v; err %v", url, err)
		}
		if url != expected {
			t.Errorf("getUrl = %v; expect %v", url, expected)
		}
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func initTestClient(fn RoundTripFunc) {
	client.client = &http.Client{
		Transport: fn,
	}
}

func authSuccessRequest(t *testing.T) RoundTripFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "http://localhost:8000/users/signin")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: io.NopCloser(strings.NewReader("{\"access_token\":\"testToken\"}")),
			// Must be set to non-nil value, or it panics
			//Header: make(http.Header),
		}
	}
}

func authFailedRequest(t *testing.T) RoundTripFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "http://localhost:8000/users/signin")
		return &http.Response{
			StatusCode: 401,
			// Send response to be tested
			Body: io.NopCloser(strings.NewReader("{\"error\":\"unauthenticated\"}")),
			// Must be set to non-nil value, or it panics
			//Header: make(http.Header),
		}
	}
}

func TestApiClient_Authenticate(t *testing.T) {
	viper.Set("ACCOUNTS_URL", "http://localhost:8000")

	initTestClient(authSuccessRequest(t))

	values := map[string]string{"email": "aa@a.a", "password": "12345678"}
	jsonData, err := json.Marshal(values)

	err = client.Authenticate(jsonData)
	if err != nil {
		t.Error(err)
	}

	// Test cache
	cache.Init()

	err = client.Authenticate(jsonData)
	if err != nil {
		t.Error(err)
	}

	token := cache.Instance.Get("token")
	if token.Value() == "" || token.Value() != "testToken" || token.IsExpired() {
		t.Error("cache failed to store token")
	}

	// Test auth failure
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expect authentication panic")
		}
	}()

	initTestClient(authFailedRequest(t))
	err = client.Authenticate(jsonData)
}
