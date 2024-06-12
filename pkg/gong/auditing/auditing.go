package auditing

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	iso8601Format = "2006-01-02T15:04:05Z"
)

func GetAuditLogs(accessKey string, secretKey string, lookupHours int64) ([]map[string]string, error) {
	combinedLogs := make([]map[string]string, 0)

	for logType := range logTypeStructMap {
		logs, err := getAuditLogsForType(accessKey, secretKey, logType, lookupHours)
		if err != nil {
			return nil, err
		}
		combinedLogs = append(combinedLogs, logs...)
	}

	return combinedLogs, nil
}

func getAuditLogsForType(accessKey string, secretKey string, logType string, lookupHours int64) ([]map[string]string, error) {
	// Calculate fromDateTime based on the lookup period in UTC time
	now := time.Now()
	fromDateTime := now.Add(-time.Duration(lookupHours) * time.Hour).Format(iso8601Format)

	url := fmt.Sprintf("https://api.gong.io/v2/logs?logType=%s&fromDateTime=%s", logType, fromDateTime)

	client := &http.Client{
		Timeout: time.Second * 50,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Setting up the basic authentication
	req.SetBasicAuth(accessKey, secretKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send HTTP request: %v\n", err) // Print the error message
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read and parse the response body to extract the error message
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		var errorResponse struct {
			RequestID string   `json:"requestId"`
			Errors    []string `json:"errors"`
		}
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error response: %v", err)
		}

		// Check each error message for the specific error
		for _, errMsg := range errorResponse.Errors {
			if strings.Contains(errMsg, "No log records found corresponding to the provided log type and time range") {
				logrus.Warn("No log records found corresponding to the provided log type and time range")
				return []map[string]string{}, nil // Return empty logs with no error
			} else {
				return nil, fmt.Errorf("failed to fetch audit logs for %s: %s", logType, errMsg)
			}
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var response struct {
		LogEntries []map[string]interface{} `json:"logEntries"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	TimeGenerated := time.Now().UTC().Format(iso8601Format)
	mappedLogs := make([]map[string]string, len(response.LogEntries))

	for i, entry := range response.LogEntries {
		logRecordMap := make(map[string]string)
		logRecordMap["TimeGenerated"] = TimeGenerated
		logRecordMap["logType"] = logType

		// Marshal the log record to JSON
		logEntryJSON, err := json.Marshal(entry)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal log entry to JSON: %v", err)
		}
		logRecordMap["logEntry"] = string(logEntryJSON)

		mappedLogs[i] = logRecordMap
	}

	return mappedLogs, nil
}
