package auditing

import "time"

type UserActivityLog struct {
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
			TableChanges []struct {
				PreSnapshotTimestamp  time.Time `json:"preSnapshotTimestamp"`
				PostSnapshotTimestamp time.Time `json:"postSnapshotTimestamp"`
				RowChanges            []struct {
					PrimaryKeyColumns []struct {
						ColumnValue string `json:"columnValue"`
						ColumnName  string `json:"columnName"`
					} `json:"primaryKeyColumns"`
					ColumnChanges []struct {
						NewValue   string `json:"newValue"`
						OldValue   any    `json:"oldValue"`
						Operation  string `json:"operation"`
						ColumnName string `json:"columnName"`
					} `json:"columnChanges"`
				} `json:"rowChanges"`
				TableName string `json:"tableName"`
			} `json:"tableChanges"`
			Action      any `json:"action"`
			HTTPRequest struct {
				ReferrerURI string `json:"referrerUri"`
				ClientIP    string `json:"clientIp"`
				Verb        string `json:"verb"`
				EndpointURI string `json:"endpointUri"`
				Body        string `json:"body"`
				Parameters  []any  `json:"parameters"`
			} `json:"httpRequest"`
			CustomData  []any `json:"customData"`
			WorkspaceID any   `json:"workspaceId"`
		} `json:"logRecord"`
		ImpersonatorUserID       string `json:"impersonatorUserId,omitempty"`
		ImpersonatorEmailAddress string `json:"impersonatorEmailAddress,omitempty"`
		ImpersonatorFullName     string `json:"impersonatorFullName,omitempty"`
		ImpersonatorCompanyID    string `json:"impersonatorCompanyId,omitempty"`
	} `json:"logEntries"`
}
