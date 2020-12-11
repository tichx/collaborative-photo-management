package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
	"github.com/my/repo/servers/gateway/handlers"
	"github.com/my/repo/servers/gateway/models/users"
	"github.com/my/repo/servers/gateway/sessions"
)

//Director is a middleware
type Director func(r *http.Request)

//main is the main entry point for the server
func main() {
	/*add code to do the following
	- Read the ADDR environment variable to get the address
	  the server should listen on. If empty, default to ":80"
	- Create a new mux for the web server.
	- Tell the mux to call your handlers.SummaryHandler function
	  when the "/v1/summary" URL path is requested.
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/

	// 1.Read the ADDR environment variable to get the address
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	TLSCERT := os.Getenv("TLSCERT")
	if len(TLSCERT) == 0 {
		log.Println("No TLSCERT environment variable found")
		os.Exit(1)
	}

	TLSKEY := os.Getenv("TLSKEY")
	if len(TLSKEY) == 0 {
		log.Println("No TLSKEY environment variable found")
		os.Exit(1)
	}

	sessionkey := os.Getenv("SESSIONKEY")
	if len(sessionkey) == 0 {
		log.Println("SESSIONKEY env variable was not set")
		os.Exit(1)
	}

	redisaddr := os.Getenv("REDISADDR")
	if len(redisaddr) == 0 {
		log.Println("REDISADDR env variable was not set")
		os.Exit(1)
		//redisaddr = "172.19.0.2:6379"
		//redisaddr = "redisServer:6379"
	}

	//3306
	dsn := os.Getenv("DSN")
	if len(dsn) == 0 {
		log.Println("DSN env variable was not set")
		os.Exit(1)
	}

	//Create DB object from SQL DB
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error opening the database: %v", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error opening the database: %v", err)
		os.Exit(1)
	}

	//When comeplete, close the db
	defer db.Close()

	msgAddrs := strings.Split(os.Getenv("MESSAGESADDR"), ",")
	if len(msgAddrs) == 0 {
		msgAddrs = append(msgAddrs, "http://micro-messaging:4000")
	}

	//Createmysqlstore
	usersStore := users.NewMySQLStore(db)

	//Create redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})

	//Create redisstore
	sessionStore := sessions.NewRedisStore(redisClient, time.Hour)

	//Create context
	context := handlers.NewContext(sessionkey, sessionStore, usersStore)

	// 2.Create a new mux for the web server.
	mux := http.NewServeMux()

	var msgServerAddrs []*url.URL
	for _, msgAddr := range msgAddrs {
		msgSerAddr, _ := url.Parse(msgAddr)
		msgServerAddrs = append(msgServerAddrs, msgSerAddr)
	}

	msgProxy := &httputil.ReverseProxy{Director: CustomDirector(msgServerAddrs, context)}
	mux.Handle("/v1/channels", msgProxy)
	mux.Handle("/v1/channels/", msgProxy)
	mux.Handle("/v1/messages", msgProxy)
	mux.Handle("/v1/messages/", msgProxy)

	// 3.Tell the mux to call your handlers.SummaryHandler function
	//mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/user/{UserID | me}")
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)

	//   4.Start a web server listening on the address you read from
	//   the environment variable, using the mux you created as
	//   the root handler. Use log.Fatal() to report any errors
	//   that occur when trying to start the web server.
	wrappedMux := handlers.NewCors(mux)

	log.Printf("server is listening at http://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, TLSCERT, TLSKEY, wrappedMux))
}

//CustomDirector takes in session context and do authentication
func CustomDirector(targets []*url.URL, context *handlers.SessionContext) Director {
	var counter int32
	counter = 0

	return func(r *http.Request) {
		targ := targets[int(counter)%len(targets)]
		atomic.AddInt32(&counter, 1)

		//Authenticate user
		sessionState := &handlers.SessionState{}
		sessions.GetState(r, context.Key, context.SessionsStore, sessionState)
		//Get user from session state
		user := sessionState.User

		//if there's a currently-authenticated user
		if strconv.Itoa(int(user.ID)) != "" {
			encoded, _ := json.Marshal(user)
			encodedUser := base64.StdEncoding.EncodeToString(encoded)
			r.Header.Add("X-User", encodedUser)
		} else {
			//remove possible X-User header if no currently-authenticated user. Prevent spoofing from client
			r.Header.Del("X-User")
		}

		// log.Println("target host: ", targ.Host)
		// log.Println("scheme: ", targ.Scheme)
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme

	}
}
