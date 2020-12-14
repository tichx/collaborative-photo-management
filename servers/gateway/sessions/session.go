package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	//- create a new SessionID
	//- save the sessionState to the store
	//- add a header to the ResponseWriter that looks like this:
	//    "Authorization: Bearer <sessionID>"
	//  where "<sessionID>" is replaced with the newly-created SessionID
	//  (note the constants declared for you above, which will help you avoid typos)
	sid, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	data := store.Save(sid, sessionState)
	if data != nil {
		return InvalidSessionID, err
	}

	w.Header().Add(headerAuthorization, schemeBearer+sid.String())
	return sid, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.
	var sidraw string
	var sid string
	var err error

	//if Authorization header exists, extract the header
	if _, exists := r.Header[headerAuthorization]; exists {
		sidraw = r.Header.Get(headerAuthorization)
		sid, err = ValidateBearerHelper(sidraw)
		if nil != err {
			return InvalidSessionID, err
		}
		//else attempt to extract from query
	} else {
		sidraw = r.URL.Query().Get("auth")
		sid, err = ValidateBearerHelper(sidraw)
		if nil != err {
			return InvalidSessionID, err
		}
	}
	validSid, err := ValidateID(sid, signingKey)
	if nil != err {
		return InvalidSessionID, err
	}
	return validSid, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//get the SessionID from the request, and get the data
	//associated with that SessionID from the store.
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	resp := store.Get(sid, sessionState)
	if resp != nil {
		return InvalidSessionID, resp
	}
	return sid, nil
}

//ValidateBearerHelper takes in a raw sid string, check if it's empty or has wrong prefix
//return sid if its in correct format, or else return ErrInvalidScheme error
func ValidateBearerHelper(sidraw string) (string, error) {
	if len(sidraw) == 0 || !(strings.HasPrefix(sidraw, schemeBearer)) {
		return "", ErrInvalidScheme
	}
	sid := strings.TrimPrefix(sidraw, schemeBearer)
	return sid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//get the SessionID from the request, and delete the
	//data associated with it in the store.
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	err = store.Delete(sid)
	if err != nil {
		return InvalidSessionID, err
	}
	return sid, nil
}
