package handlers

import (
	"assignments-tichx/servers/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

// SessionState is capable of tracking the the
// time at which this session began, and the user
type SessionState struct {
	Time time.Time   `json:"time"`
	User *users.User `json:"user"`
}
