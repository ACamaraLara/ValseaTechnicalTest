package tester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "http://localhost:8080"

func printResponseJson(responseObject interface{}) {
	// Indent Json to show the user pretty-printed Json.
	prettyJSON, _ := json.MarshalIndent(responseObject, "", "  ")
	fmt.Printf("%+v\n\n", string(prettyJSON))
}

func printResponseBody(body io.ReadCloser) {
	data, err := io.ReadAll(body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Printf("Response Body: %s\n\n", string(data))
}

func getAccountByOwner(message string, accountRepo AccountsRepo) (*Account, string) {
	fmt.Print(message)
	owner := InputString()

	storedAccount, exists := accountRepo[owner]
	if !exists {
		return nil, ""
	}
	return storedAccount, owner
}

func makeRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	url := baseURL + endpoint
	var body []byte
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return &http.Response{}, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = jsonPayload
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return &http.Response{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &http.Response{}, fmt.Errorf("failed to execute request: %w", err)
	}
	return resp, nil
}
