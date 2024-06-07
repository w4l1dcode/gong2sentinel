package auditing

import "time"

type ExternallySharedCallPlay struct {
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
			TimeBasedSecureSharingID string    `json:"time_based_secure_sharing_id"`
			CallID                   string    `json:"call_id"`
			VideoPlayerInstanceID    string    `json:"video_player_instance_id"`
			SequenceNum              string    `json:"sequence_num"`
			PlaySpeed                float64   `json:"play_speed"`
			Device                   string    `json:"device"`
			StartTime                float64   `json:"start_time"`
			EndTime                  float64   `json:"end_time"`
			EventTimeOnDevice        time.Time `json:"event_time_on_device"`
			Offline                  bool      `json:"offline"`
			Live                     bool      `json:"live"`
		} `json:"logRecord"`
		UserFullName string `json:"userFullName,omitempty"`
	} `json:"logEntries"`
}
