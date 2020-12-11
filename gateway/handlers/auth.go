package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/my/repo/servers/gateway/models/users"
	"github.com/my/repo/servers/gateway/sessions"
)

//define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

// UsersHandler handles requests for the "users" resource
func (context *SessionContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-type")

		//If the request's Content-Type header does not start with application/json,
		//respond with status code 415,
		//and a message indicating that the request body must be in JSON.
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("The request body must be in JSON format"))
			return
		}
		//decode new user struct from request body
		decoder := json.NewDecoder(r.Body)

		var nu users.NewUser
		err := decoder.Decode(&nu)
		if err != nil {
			//respond with status code 400
			http.Error(w, "Failed to decode", http.StatusBadRequest)
		} else {
			//converts new user struct to user struct
			user, err := nu.ToUser()
			if err != nil {
				log.Printf("could not convert the user: %v", err.Error())
				//respond with status code 400
				http.Error(w, "could not convert the user", http.StatusBadRequest)
				return
			}
			//Insert decoded user into db
			uResp, insertErr := context.UsersStore.Insert(user)
			if insertErr != nil {
				log.Printf("Insert Failed. Error was: %v", insertErr.Error())
				http.Error(w, "Error Inserting into Database", http.StatusInternalServerError)
				return
			}

			//Generate and store new session based on the session context
			newSessState := SessionState{time.Now(), *uResp}
			_, sessionErr := sessions.BeginSession(context.Key, context.SessionsStore, newSessState, w)
			if sessionErr != nil {
				//delete the inserted user if this function doesn't run correctly
				context.UsersStore.Delete(uResp.ID)
				//if unauthorized, set status code 401
				http.Error(w, "Error Beginning a new session: "+sessionErr.Error(), http.StatusUnauthorized)

				return
			}
			w.Header().Set("Content-Type", "application/json")
			//write status 201: status created
			w.WriteHeader(http.StatusCreated)
			//Encode user. HashPass and Email already defined as hidden in User struct
			json.NewEncoder(w).Encode(uResp)
		}

	} else {
		//invalid HTTP request, set status code 405
		http.Error(w, "Only POST is supported for this handler", http.StatusMethodNotAllowed)
	}
}

//SpecificUserHandler handles request for a specific user
func (context *SessionContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	currentSessionState := &SessionState{}
	_, err := sessions.GetState(r, context.Key, context.SessionsStore, currentSessionState)
	if err != nil {
		//if unauthorized, set status code 401
		http.Error(w, "Unauthorized session: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		//Takes either "me" or a number, Get the requested ID
		requestedID, err := RequestedIDHelper(r, currentSessionState.User.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		//Look for the user
		print("here's the requested ID: ", requestedID)
		requestedUser, GetByIDErr := context.UsersStore.GetByID(requestedID)
		if GetByIDErr != nil {
			print("error initiated")
			//set status code 404
			http.Error(w, "404: User not found", http.StatusNotFound)
			return
		}
		//Return the user in JSON
		w.Header().Set("Content-Type", "application/json")
		//set status code 200
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(requestedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//Process updates from authenicated user on objects and applies it to a stored user
	if r.Method == http.MethodPatch {
		//If the request's Content-Type header does not start with application/json,
		//respond with status code http.StatusUnsupportedMediaType (415),
		//and a message indicating that the request body must be in JSON.
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "The request body must be in JSON!", http.StatusUnsupportedMediaType)
			return
		}

		//Get the requested ID
		requestedID, err := RequestedIDHelper(r, currentSessionState.User.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		//check if cliamed user is this session's user
		if currentSessionState.User.ID == requestedID {

			//Decode the given updates from the request
			var givenUpdates users.Updates
			err := json.NewDecoder(r.Body).Decode(&givenUpdates)
			if err != nil {
				http.Error(w, "can't decode the update: "+err.Error(), http.StatusInternalServerError)
				return
			}
			//update current user struct
			errUpdate := currentSessionState.User.ApplyUpdates(&givenUpdates)
			if errUpdate != nil {
				//set status code 403
				http.Error(w, "error when updating profile. Check First&Last Name", http.StatusForbidden)
				return
			}

			//update user in database
			updatedUser, err := context.UsersStore.Update(requestedID, &givenUpdates)
			if err != nil {
				print("encountered error during update: ", err)
				http.Error(w, "encountered error during update", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			//set status code 200
			w.WriteHeader(http.StatusOK)
			ecErr := json.NewEncoder(w).Encode(updatedUser)
			if ecErr != nil {
				http.Error(w, "Error encoding", http.StatusNotFound)
				return
			}
			//claimed user is not the actual user in this session
		} else {
			http.Error(w, "You are not authorized to make this change", http.StatusForbidden)
		}
	}
	w.WriteHeader(http.StatusMethodNotAllowed)

}

//SessionsHandler handles requests for the "sessions" resource,
//and allows clients to begin a new session using an existing user's credentials.
func (context *SessionContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		//if wrong request body type, set status code 415
		if contentType != "application/json" {
			http.Error(w, "The request body must be in JSON!", http.StatusUnsupportedMediaType)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var cred users.Credentials
		err := decoder.Decode(&cred)
		if err != nil {
			//respond with status code 400
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)

			return
		}

		//Get user struct with email from credential
		user, userError := context.UsersStore.GetByEmail(cred.Email)
		if userError != nil {
			// sleep 1 second to keep error response time consistent. Avoid hacking by error
			time.Sleep(1 * time.Second)
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		//Authenticate user with pwd in credential
		authErr := user.Authenticate(cred.Password)
		if authErr != nil {
			// sleep 1 second to keep error response time consistent. Avoid hacking by error
			time.Sleep(1 * time.Second)
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		//Create a new session state based on the current time
		newSessionState := SessionState{time.Now(), *user}
		//Begin a new session based on the session context
		_, sessionErr := sessions.BeginSession(context.Key, context.SessionsStore, newSessionState, w)
		if sessionErr != nil {
			//set status code 500
			http.Error(w, "Error Beginning a new session", http.StatusUnauthorized)
			return
		}

		//Log Successful sign-in
		//by default takes in RemoteAddr
		//if X-Forwarded-For header is included, use the first IP address in the list
		ipAddr := r.RemoteAddr
		if r.Header.Get("X-Forwarded-For") != "" {
			ipAddr = strings.Split(r.Header.Get("X-Forwarded-For"), ", ")[0]
		}
		_, err = context.UsersStore.InsertSignIn(user.ID, ipAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//respond to client with encoded user
		w.Header().Set("Content-Type", "application/json")
		//set status code 201
		w.WriteHeader(http.StatusCreated)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		//if not POST,set status code 405
		http.Error(w, "Invalid HTTP method. Only GET supported", http.StatusMethodNotAllowed)
	}

	return
}

//SpecificSessionHandler handles requests related to a specific authenticated session
func (context *SessionContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	currentSessionState := &SessionState{}
	_, err := sessions.GetState(r, context.Key, context.SessionsStore, currentSessionState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodDelete {
		_, identifier := filepath.Split(r.URL.Path)
		if identifier != "mine" {
			//set status code 403
			http.Error(w, "Forbidden status", http.StatusForbidden)
			return
		}
		_, err := sessions.EndSession(r, context.Key, context.SessionsStore)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//Confimation message
		w.Write([]byte("signed out"))
	}
	w.WriteHeader(http.StatusMethodNotAllowed)

}

//RequestedIDHelper returns an int64 of the requested userID from the URL path
//for internal use only
func RequestedIDHelper(r *http.Request, authenticatedUserID int64) (int64, error) {
	//UserID is the user param passed in by client
	//filepath.Split returns path and file (ending)
	//doc: https://golang.org/pkg/path/filepath/#Split
	_, UserID := filepath.Split(r.URL.Path)

	var finalReqestedID int64
	//Could be using me or an actual number
	if UserID == "me" {
		//if it is "me" grab the ID of the currently authenticated user
		finalReqestedID = authenticatedUserID
	} else {
		//doc: https://golang.org/pkg/strconv/
		finalReqestedIDInt, err := strconv.ParseInt(UserID, 10, 64)
		if err != nil {
			return 0, err
		}

		finalReqestedID = int64(finalReqestedIDInt)
	}
	return finalReqestedID, nil
}
