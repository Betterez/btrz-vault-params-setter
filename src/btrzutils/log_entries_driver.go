package btrzutils

// LogEntriesConnection - represent a log entry api connection
type LogEntriesConnection struct {
	apiKey    string
	accountId string
}

// CreateConnection - returns new connection or an error
func CreateConnection(APIKey, accountId string) (*LogEntriesConnection, error) {
	result := &LogEntriesConnection{
		apiKey:    APIKey,
		accountId: accountId,
	}

	return result, nil
}
