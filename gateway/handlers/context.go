package handlers

import (
	"github.com/my/repo/servers/gateway/models/users"
	"github.com/my/repo/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

//SessionContext captures the signing key, session and user info
type SessionContext struct {
	Key           string               `json:"-"`
	SessionsStore *sessions.RedisStore `json:"session"`
	UsersStore    *users.MySQLStore    `json:"user"`
}

//NewContext creates a new context if given a key, sessionstore and userstore
func NewContext(key string, ss *sessions.RedisStore, us *users.MySQLStore) *SessionContext {
	if ss == nil || us == nil || key == "" {
		return nil
	}

	context := SessionContext{key, ss, us}
	return &context
}
