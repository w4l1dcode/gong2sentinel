package calls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	iso8601Format = "2006-01-02T15:04:05Z"
)

// PostRequestBody Define the struct for the POST request body
type PostRequestBody struct {
	Filter struct {
		CallIds []string `json:"callIds"`
	} `json:"filter"`
}

// ResponseBody represents the response structure of the POST request
type ResponseBody struct {
	RequestID      string              `json:"requestId"`
	CallAccessList []map[string]string `json:"callAccessList"`
}

// GetUserAccess Function to make a POST request with filtered call IDs
func GetUserAccess(accessKey string, secretKey string, callIds []string) ([]map[string]string, error) {
	now := time.Now().UTC().Format(iso8601Format)
	url := "https://api.gong.io/v2/calls/users-access"
	client := &http.Client{}

	// Create the POST request body with the filtered call IDs
	postRequestBody := &PostRequestBody{
		Filter: struct {
			CallIds []string `json:"callIds"`
		}{CallIds: callIds},
	}

	// Convert request body to JSON
	jsonData, err := json.Marshal(postRequestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body to JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set basic authentication
	req.SetBasicAuth(accessKey, secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to send POST request: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal the response body into the structured representation
	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// Convert CallAccessList to the desired format
	callAccessList := make([]map[string]string, len(responseBody.CallAccessList))
	for i, item := range responseBody.CallAccessList {
		callAccessList[i] = map[string]string{
			"TimeGenerated":  now,
			"requestId":      item["requestId"],
			"callAccessList": item["callAccessList"],
		}
	}
	return callAccessList, nil
}
