package btrzutils

// LogsEntriesLog - helper
type LogsEntriesLog struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Links  []LogLink `json:"links"`
	Tokens []string  `json:"tokens"`
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

func (lel LogsEntriesLog) String() string {
	return lel.ID + "-" + lel.Name
}

// HasTokens does this log has a token
func (lel *LogsEntriesLog) HasTokens() bool {
	return len(lel.Tokens) > 0
}
