package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
)

func ProcessPayment(userName string, amount float32, response interface{}) error {
	url := "https://dummy-payment-server.herokuapp.com/payment"
	paymentType := "DEBIT"
	if amount > 0 {
		paymentType = "CREDIT"
	}
	type Payload struct {
		UserName    string  `json:"user_name"`
		PaymentType string  `json:"payment_type"`
		Amount      float64 `json:"amount"`
	}

	payload := Payload{
		UserName:    userName,
		PaymentType: paymentType,
		Amount:      math.Abs(float64(amount)),
	}

	log.Println("payload", payload)

	req, err := newRequest(http.MethodPost, url, payload)

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("Error in request")
	}

	if response == nil {
		return nil
	}

	if w, ok := response.(io.Writer); ok {
		io.Copy(w, resp.Body)
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(response)
}

// NewRequest constructs a request
// Convert payload to a JSON
func newRequest(method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	return http.NewRequest(method, url, buf)
}
