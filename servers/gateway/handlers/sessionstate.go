package handlers

import (
	"assignments-ydf1014/servers/gateway/models/users"
	"log"
	"time"
)

//session state struct for this web server
//see the assignment description for the fields included
//note that other packages can only see exported fields!

type SessionState struct {
	BeginTime time.Time   `json:"beginTime"`
	User      *users.User `json:"user"`
}

func InitializeState(usr *users.User) *SessionState {
	log.Printf("sessionstate")
	return &SessionState{
		BeginTime: time.Now(),
		User:      usr,
	}
}
