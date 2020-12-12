package handlers

import (
	"assignments-tichx/servers/gateway/models/users"
	"assignments-tichx/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

//HandlerContext holds context values
//used by multiple handler functions.
type HandlerContext struct {
	Key          string
	UserStore    users.Store
	SessionStore sessions.Store
}
