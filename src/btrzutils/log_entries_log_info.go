package btrzutils

type LogsEntriesLog struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Links []LogLink `json:"links"`
}

type LogLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
