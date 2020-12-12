package handlers

import (
	"assignments-tichx/servers/gateway/models/users"
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//AllowOrigin gives access control origin
const AllowOrigin = "Access-Control-Allow-Origin"

//AllowHeaders gives a const string
const AllowHeaders = "Access-Control-Allow-Headers"

//AllowMethods gives a const string
const AllowMethods = "Access-Control-Allow-Methods"

//AllowMethodsTypes gives a const string
const AllowMethodsTypes = "GET, PUT, POST, PATCH, DELETE"

//AccessExposeHeaders gives a const string
const AccessExposeHeaders = "Access-Control-Expose-Headers"

//ContentTypeeAuthorization gives a const string
const ContentTypeeAuthorization = "Content-Type, Authorization"

//AccessControlMaxAge gives a const string
const AccessControlMaxAge = "Access-Control-Max-Age"

//XForwarded is the name of header X-Forwarded-For
const XForwarded = "X-Forwarded-For"

// CORS handles cors requests
type CORS struct {
	Handle http.Handler
}

// ServeHTTP Adds CORS Headers to the response
func (c *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(AccessExposeHeaders, authorization)
	w.Header().Set(AllowOrigin, "*")
	w.Header().Set(AllowHeaders, ContentTypeeAuthorization)
	w.Header().Set(AllowMethods, AllowMethodsTypes)
	w.Header().Set(AccessControlMaxAge, "600")

	// preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	//insert signin
	if r.URL.Path == "/v1/sessions" && http.MethodPost == r.Method {
		if contentTypeJSON != r.Header.Get(headerContentType) {
			http.Error(w, "CORS Error: the request format incorrect", http.StatusUnsupportedMediaType)
			return
		}
		cr := &users.Credentials{}
		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(body)))
		err := json.Unmarshal(body, cr)
		if err != nil {
			http.Error(w, "Cors: JSON conversion error", http.StatusBadRequest)
			return
		}
		db, err := sql.Open("mysql", os.Getenv("DSN"))
		if err != nil {
			http.Error(w, "Cors Error: db conn failed", http.StatusInternalServerError)
		}
		defer db.Close()
		scanner, err := db.Query("select id from users where email=?", cr.Email)
		if scanner.Next() {
			if err == nil {
				uid := int64(-1)
				err = scanner.Scan(&uid)
				if err == nil {
					// use RemoteAddr
					addr := r.RemoteAddr
					// if x-forwareded exisist use that instead
					if 0 != len(r.Header.Get(XForwarded)) {
						addr = r.Header.Get(XForwarded)
					}
					// insert to user_signin
					_, err = db.Exec("insert into user_signin (id, date_time, ip_addr) values (?, ?, ?)", uid, time.Now().Format(time.RFC3339), addr)
				}
			}
		}
	}
	c.Handle.ServeHTTP(w, r)
}
