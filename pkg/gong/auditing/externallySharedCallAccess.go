package auditing

import "time"

type ExternallySharedCallAccess struct {
	RequestID string `json:"requestId"`
	Records   struct {
		TotalRecords      int `json:"totalRecords"`
		CurrentPageSize   int `json:"currentPageSize"`
		CurrentPageNumber int `json:"currentPageNumber"`
	} `json:"records"`
	LogEntries []struct {
		UserEmailAddress string    `json:"userEmailAddress"`
		EventTime        time.Time `json:"eventTime"`
		LogRecord        struct {
			CallID                   string `json:"call_id"`
			TimeBasedSecureSharingID string `json:"time_based_secure_sharing_id"`
			PageViewerIP             string `json:"page_viewer_ip"`
		} `json:"logRecord"`
		UserFullName string `json:"userFullName,omitempty"`
	} `json:"logEntries"`
}
