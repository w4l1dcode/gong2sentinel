package auditing

import "time"

type AccessLog struct {
	RequestID string `json:"requestId"`
	Records   struct {
		TotalRecords      int    `json:"totalRecords"`
		CurrentPageSize   int    `json:"currentPageSize"`
		CurrentPageNumber int    `json:"currentPageNumber"`
		Cursor            string `json:"cursor"`
	} `json:"records"`
	LogEntries []struct {
		UserID           string    `json:"userId"`
		UserEmailAddress string    `json:"userEmailAddress"`
		UserFullName     string    `json:"userFullName"`
		EventTime        time.Time `json:"eventTime"`
		LogRecord        struct {
			ResponseHeaders struct {
				XTraceid    string `json:"x-traceid"`
				ContentType string `json:"content-type"`
				XIid        string `json:"x-iid"`
			} `json:"response_headers"`
			Protocol       string `json:"protocol"`
			Method         string `json:"method"`
			RequestHeaders struct {
				Referer       string `json:"referer"`
				XForwardedFor string `json:"x-forwarded-for"`
				UserAgent     string `json:"user-agent"`
			} `json:"request_headers"`
			ElapsedTime  int    `json:"elapsed_time"`
			RequestedURL string `json:"requested_url"`
			Message      string `json:"message"`
			Mdc          struct {
				Xtid string `json:"xtid"`
			} `json:"mdc"`
			ContentLength int    `json:"content_length"`
			RequestedURI  string `json:"requested_uri"`
			Status        int    `json:"status"`
		} `json:"logRecord"`
	} `json:"logEntries"`
}
