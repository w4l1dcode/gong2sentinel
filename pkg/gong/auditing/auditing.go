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

func GetAuditLogsForType(accessKey string, secretKey string, logType string, lookupHours int64) ([]map[string]string, error) {
	now := time.Now().UTC()
	fromDateTime := now.Add(-time.Duration(lookupHours) * time.Hour).Format(iso8601Format)

	url := fmt.Sprintf("https://api.gong.io/v2/logs?logType=%s&fromDateTime=%s", logType, fromDateTime)
	logrus.Infof("Fetching URL for logType %s: %s", logType, url)

	client := &http.Client{
		Timeout: time.Second * 50,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.SetBasicAuth(accessKey, secretKey)

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Failed to send HTTP request for logType %s: %v", logType, err)
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

		for _, errMsg := range errorResponse.Errors {
			if strings.Contains(errMsg, "No log records found corresponding to the provided log type and time range") {
				logrus.Warnf("No log records found for logType %s", logType)
				return []map[string]string{}, nil
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

		logEntryJSON, err := json.Marshal(entry)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal log entry to JSON: %v", err)
		}
		logRecordMap["logEntry"] = string(logEntryJSON)

		mappedLogs[i] = logRecordMap
	}

	return mappedLogs, nil
}
