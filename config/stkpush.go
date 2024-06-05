package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RequestData struct {
	PhoneNumber int `json:"phone_number"`
	Amount      int `json:"amount"`
}

func (c *Config) GetSTKPUSH(w http.ResponseWriter, r *http.Request) {
	accessToken, err := c.getAuth()
	if err != nil {
		http.Error(w, "Error getting access token", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData RequestData
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	phoneNumber := requestData.PhoneNumber
	amount := requestData.Amount

	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"

	businessShortCode := 174379
	partyA := phoneNumber
	partyB := 174379

	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)

	password := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d%s%s", businessShortCode, c.getPassKey(), timestamp)))

	payload := strings.NewReader(fmt.Sprintf(`{
		"BusinessShortCode": %d,
		"Password": "%s",
		"Timestamp": "%s",
		"TransactionType": "CustomerPayBillOnline",
		"Amount": %d,
		"PartyA": %d,
		"PartyB": %d,
		"PhoneNumber": %d,
		"CallBackURL": "",
		"AccountReference": "CompanyXLTD",
		"TransactionDesc": "Payment of Boogs"
	}`, businessShortCode, password, timestamp, amount, partyA, partyB, phoneNumber))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
