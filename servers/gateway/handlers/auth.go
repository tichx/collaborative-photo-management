package handlers

import (
	"assignments-ydf1014/servers/gateway/models/users"
	"assignments-ydf1014/servers/gateway/sessions"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

//these handler is HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.
func (context *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("users handler.")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("request must be post request"))
		return
	}
	contentTypeHeader := r.Header.Get("Content-type")
	contentType := strings.Split(contentTypeHeader, ",")[0]
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("request must be in JSON"))
		return
	}
	jsonDecoder := json.NewDecoder(r.Body)
	nuUser := &users.NewUser{}
	err := jsonDecoder.Decode(nuUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot decode request"))
		return
	}

	err = nuUser.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot validate user"))
		return
	}
	user, err := nuUser.ToUser()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot conver to user"))
		return
	}

	user, err = context.User.Insert(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot insert user to database"))
		return
	}

	newState := InitializeState(user)
	_, sessionErr := sessions.BeginSession(context.Key, context.Session, *newState, w)
	if sessionErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot begin new session"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot encode user"))
		return
	}

}

func (context *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	sessionState := &SessionState{}
	_, err := sessions.GetSessionID(r, context.Key)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User is unauthorized"))
		return
	}

	pathBase := path.Base(r.URL.Path)
	if r.Method == http.MethodGet {
		_, err = sessions.GetState(r, context.Key, context.Session, sessionState)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user state not found"))
			return
		}

		userID := sessionState.User.ID
		if pathBase != "me" {
			uid, err := strconv.Atoi(pathBase)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to convert userid"))
				return
			}
			userID = int64(uid)
		}

		user, err := context.User.GetByID(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("failed to find user by id"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
		return
	} else if r.Method == http.MethodPatch {
		_, err = sessions.GetState(r, context.Key, context.Session, sessionState)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user state not found"))
			return
		}

		uid, _ := strconv.ParseInt(pathBase, 10, 64)
		user := sessionState.User
		userID := user.ID
		if pathBase != "me" && uid != userID {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("user id incorrect"))
			return
		}
		if r.Header.Get("Content-type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("content type must be json"))
			return
		}
		update := &users.Updates{}
		err := json.NewDecoder(r.Body).Decode(update)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("cannot update json"))
			return
		}
		err = user.ApplyUpdates(update)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid name"))
			return
		}
		user, err = context.User.Update(userID, update)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("fail to update user"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("cannot encode user"))
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("must be get or patch method"))
		return
	}
}

func (context *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("request body must be in json"))
			return
		}
		decoder := json.NewDecoder(r.Body)
		cred := &users.Credentials{}
		err := decoder.Decode(cred)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("cannot decode json"))
			return
		}
		user, err := context.User.GetByEmail(cred.Email)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user unauthorized"))
			return
		}
		err = user.Authenticate(cred.Password)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user unauthenticated"))
			return
		}
		_, err = sessions.BeginSession(context.Key, context.Session, InitializeState(user), w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("fail to begin session"))
			return
		}
		ipAddr := r.Header.Get("X-Forwarded-For")
		currentAddr := r.RemoteAddr
		if len(ipAddr) != 0 {
			currentAddr = ipAddr
		}
		userID := strconv.FormatInt(user.ID, 10)
		logSign := users.SignInLog{
			ID:       userID,
			DateTime: time.Now().String(),
			IPAddr:   currentAddr,
		}

		context.User.LogSignin(&logSign)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("fail encode user"))
			return
		}

	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (context *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		lastSegment := path.Base(r.URL.Path)
		if lastSegment != "mine" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("last path segment isn't mine"))
			return
		}
		_, err := sessions.EndSession(r, context.Key, context.Session)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("fail to close session"))
			return
		}
		w.Write([]byte("signed out"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
