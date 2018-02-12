package btrzutils

// LogsEntriesLog - helper
type LogsEntriesLog struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Links       []LogLink           `json:"links"`
	Tokens      []string            `json:"tokens"`
	LogSetsInfo []LogsEntriesLogSet `json:"logsets_info"`
}

func (le LogsEntriesLog) String() string {
	result := le.Name
	if len(le.LogSetsInfo) > 0 {
		result += " [" + le.LogSetsInfo[0].Name + "]"
	}
	if len(le.Tokens) > 0 {
		result += " " + le.Tokens[0]
	}
	return result
}

// LogEntriesLogsResponse -helper container
type LogEntriesLogsResponse struct {
	Logs []LogsEntriesLog `json:"logs"`
}

// LogLink - helper
type LogLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// LogsEntriesLogResponse - helper
type LogsEntriesLogResponse struct {
	Log LogsEntriesLog `json:"log"`
}

// HasTokens  - does this log has a token
func (le *LogsEntriesLog) HasTokens() bool {
	return len(le.Tokens) > 0
}
