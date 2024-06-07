package calls

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetCallIDs function to make a GET request and retrieve call IDs
func GetCallIDs(accessKey string, secretKey string) ([]string, error) {
	url := "https://api.gong.io/v2/calls"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set basic authentication
	req.SetBasicAuth(accessKey, secretKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch call IDs: %s", resp.Status)
	}

	var callIDs []string
	var response struct {
		Calls []struct {
			ID string `json:"id"`
		} `json:"calls"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	for _, call := range response.Calls {
		callIDs = append(callIDs, call.ID)
	}

	return callIDs, nil
}
