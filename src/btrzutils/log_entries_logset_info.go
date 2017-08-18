package btrzutils

// LogsEntriesLogSet  log set info
type LogsEntriesLogSet struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	LogsInfo    []LogsEntriesLog `json:"logs_info"`
}

func (ls LogsEntriesLogSet) String() string {
	return ls.ID + " " + ls.Name
}

// LogEntriesLogSetResponse - response helper
type LogEntriesLogSetResponse struct {
	Logsets []LogsEntriesLogSet `json:"logsets"`
}
