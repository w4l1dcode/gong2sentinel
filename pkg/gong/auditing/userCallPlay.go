package auditing

import "time"

type UserCallPlay struct {
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
			CallID                string    `json:"call_id"`
			VideoPlayerInstanceID string    `json:"video_player_instance_id"`
			SequenceNum           int       `json:"sequence_num"`
			PlaySpeed             float64   `json:"play_speed"`
			Device                string    `json:"device"`
			StartTime             float64   `json:"start_time"`
			EndTime               float64   `json:"end_time"`
			EventTimeOnDevice     time.Time `json:"event_time_on_device"`
			Offline               bool      `json:"offline"`
			Live                  bool      `json:"live"`
		} `json:"logRecord"`
	} `json:"logEntries"`
}
