package handlers

import (
	"assignments-tichx/servers/gateway/models/users"
	"assignments-tichx/servers/gateway/sessions"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrigin = "*"
const contentTypeJSON = "application/json"
const headerContentType = "Content-Type"
const authorization = "Authorization"

// UsersHandler handles new user registration
func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if http.MethodPost != r.Method {
		http.Error(w, "Method Must be post", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get(headerContentType) != contentTypeJSON {
		http.Error(w, "Response is not formatted in json", http.StatusUnsupportedMediaType)
		return
	}
	acc := &users.NewUser{}
	err := json.NewDecoder(r.Body).Decode(acc)
	if err != nil {
		http.Error(w, "JSON post decode error", http.StatusBadRequest)
		return
	}
	err = acc.Validate()
	if err != nil {
		http.Error(w, "validation went wrong", http.StatusBadRequest)
		return
	}
	user, err := acc.ToUser()
	if err != nil {
		http.Error(w, "ToUser went wrong", http.StatusInternalServerError)
		return
	}
	nu, err := ctx.UserStore.Insert(user)
	if err != nil {
		http.Error(w, "Insertion went wrong", http.StatusBadRequest)
		return
	}
	sess := &SessionState{
		Time: time.Now(),
		User: nu,
	}
	_, err = sessions.BeginSession(ctx.Key, ctx.SessionStore, sess, w)
	if err != nil {
		http.Error(w, "Session Begin went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(nu)
	if err != nil {
		http.Error(w, "JSON encode went wrong", http.StatusInternalServerError)
		return
	}
}

// SpecificUserHandler handles speecific user get/patch actions
func (ctx *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	path := path.Base(r.URL.Path)
	sid := sessions.SessionID(strings.TrimPrefix(r.Header.Get(authorization), "Bearer "))
	sess := &SessionState{}
	err := ctx.SessionStore.Get(sid, sess)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	uid := int64(-1)
	if "me" != path {
		uid, err = strconv.ParseInt(path, 10, 64)
		if err != nil {
			http.Error(w, "String conversion error", http.StatusInternalServerError)
		}
	} else {
		uid = sess.User.ID
	}
	if http.MethodPatch == r.Method {
		if contentTypeJSON != r.Header.Get(headerContentType) {
			http.Error(w, "Request must be formatted in json", http.StatusUnsupportedMediaType)
			return
		}
		if uid != sess.User.ID {
			http.Error(w, "User does not match", http.StatusForbidden)
			return
		}
		updates := &users.Updates{}
		err := json.NewDecoder(r.Body).Decode(updates)
		if err != nil {
			http.Error(w, "Decoding body error", http.StatusBadRequest)
			return
		}
		user, err := ctx.UserStore.Update(sess.User.ID, updates)
		if err != nil {
			http.Error(w, "update user went wrong", http.StatusInternalServerError)
			return
		}
		sess.User = user
		err = ctx.SessionStore.Save(sid, sess)
		if err != nil {
			http.Error(w, "update user went wrong", http.StatusInternalServerError)
			return
		}
		w.Header().Set(headerContentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "JSON encode went wrong", http.StatusInternalServerError)
			return
		}
	}
	if r.Method == http.MethodGet {
		user, err := ctx.UserStore.GetByID(uid)
		if err != nil {
			http.Error(w, "Unable to find user", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set(headerContentType, contentTypeJSON)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "JSON encode went wrong", http.StatusInternalServerError)
			return
		}
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPatch {
		http.Error(w, "This method is forbidden", http.StatusMethodNotAllowed)
		return
	}
}

// SessionsHandler is used to create a session
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if http.MethodPost != r.Method {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	if contentTypeJSON != r.Header.Get(headerContentType) {
		http.Error(w, "Request must be in JSON", http.StatusUnsupportedMediaType)
		return
	}
	if r.Body == nil {
		http.Error(w, "Request is empty", http.StatusUnsupportedMediaType)
		return
	}

	c := &users.Credentials{}
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		http.Error(w, "Decoding body error: "+err.Error(), http.StatusBadRequest)
		return
	}
	user, err := ctx.UserStore.GetByEmail(c.Email)
	if err != nil {
		bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	err = user.Authenticate(c.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	sess := &SessionState{
		Time: time.Now(),
		User: user,
	}
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	_, err = sessions.BeginSession(ctx.Key, ctx.SessionStore, sess, w)
	if err != nil {
		http.Error(w, "Session begin went wrong", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "JSON encode went wrong", http.StatusInternalServerError)
		return
	}

}

// SpecificSessionHandler will log user out
func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	url := path.Base(r.URL.Path)
	if http.MethodDelete != r.Method {
		http.Error(w, "SpecificSessionHandler Error: Only DELETE is allowed", http.StatusMethodNotAllowed)
		return
	}
	if "mine" != url {
		http.Error(w, "SpecificSessionHandler Error: incorrect request base url", http.StatusForbidden)
		return
	}
	if http.MethodDelete == r.Method {
		_, err := sessions.EndSession(r, ctx.Key, ctx.SessionStore)
		if err != nil {
			http.Error(w, "SpecificSessionHandler Error: Session end went wrong "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("signed out"))
	}
}
