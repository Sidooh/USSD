package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Data struct {
	SessionId   string `json:"sessionId"`
	ServiceCode string `json:"serviceCode"`
	PhoneNumber string `json:"phoneNumber"`
	NetworkCode string `json:"networkCode"`
	Text        string `json:"text"`
}

func decodeData(r *http.Request) *Data {
	content := r.Header.Get("Content-Type")

	if content == "application/json" {
		decoder := json.NewDecoder(r.Body)
		var t Data
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		return &t
	} else if content == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		return &Data{
			SessionId:   r.FormValue("sessionId"),
			ServiceCode: r.FormValue("serviceCode"),
			PhoneNumber: r.FormValue("phoneNumber"),
			NetworkCode: r.FormValue("networkCode"),
			Text:        r.FormValue("text"),
		}
	}

	return nil
}

func ussdHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/ussd" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	data := decodeData(r)

	fmt.Fprintln(w, processAndRespond(data.NetworkCode, data.PhoneNumber, data.SessionId, data.Text))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	initUssd()

	fmt.Printf("Starting server at port %v\n", port)

	http.HandleFunc("/api/ussd", ussdHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
