package handlers

import (
	"assignments-ydf1014/servers/gateway/models/users"
	"assignments-ydf1014/servers/gateway/sessions"
	"log"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store
type Context struct {
	Key         string
	Session     sessions.Store
	User        users.Store
	SocketStore *SocketStore
	// Notifier:    *Notifier,
}

func Initialize(session sessions.Store, user users.Store, key string, SocketStore *SocketStore) *Context {
	log.Printf("context initialization")
	return &Context{
		Key:         key,
		Session:     session,
		User:        user,
		SocketStore: SocketStore,
		// Notifier: notifier,
	}
}
