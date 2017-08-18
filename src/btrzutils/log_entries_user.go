package btrzutils

import ()

// LogEntryUser - repesent user of le
type LogEntryUser struct {
	UserID    string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	LoginName string `json:"login_name"`
}

type usersResponse struct {
	Users []LogEntryUser `json:"users"`
}

func (user LogEntryUser) String() string {
	result := user.UserID
	if user.FirstName != "" {
		result += ", " + user.FirstName
		if user.LastName != "" {
			result += " " + user.LastName
		}
	}
	if user.Email != "" {
		result += ", " + user.Email
	}
	return result
}
