package main

import (
	"assignments-tichx/servers/gateway/models/users"
	"assignments-tichx/servers/gateway/sessions"
	"assignments-tichx/servers/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

//main is the main entry point for the server
func main() {
	const headerCORS = "Access-Control-Allow-Origin"
	const corsAnyOrigin = "*"

	addr := os.Getenv("ADDR")

	if len(addr) == 0 {
		addr = ":443"
	}
	MESSAGESADDR := strings.Split(os.Getenv("MESSAGESADDR"), ",")
	PHOTOSADDR := strings.Split(os.Getenv("PHOTOSADDR"), ",")

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	SessionKey := os.Getenv("SESSIONKEY")
	if len(SessionKey) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Session key is missing.")
		os.Exit(1)
	}
	RedisAddr := os.Getenv("REDISADDR")
	if len(RedisAddr) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Redis address is missing")
		os.Exit(1)
	}
	DSN := os.Getenv("DSN")
	if len(DSN) == 0 {
		fmt.Fprintf(os.Stderr, "Error: DSN is missing")
		os.Exit(1)
	}
	if len(tlsKeyPath) == 0 || len(tlsCertPath) == 0 {
		fmt.Fprintf(os.Stderr, "Error: both certificates path are missing, please declare as env vars.")
		os.Exit(1)
	}

	redis := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
	})
	pong, err := redis.Ping().Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Redis Connection: %v, %v", pong, err)
		os.Exit(1)
	}

	redisStore := sessions.NewRedisStore(redis, time.Hour)

	db, err := sql.Open("mysql", DSN)
	if err != nil {
		os.Stderr.WriteString("Unable to open database: " + err.Error())
		os.Exit(1)
	}
	userStore := users.NewMySQLStore(db)
	context := &handlers.HandlerContext{
		Key:          SessionKey,
		SessionStore: redisStore,
		UserStore:    userStore,
	}

	summaryAddr := os.Getenv("SUMMARYADDR")
	summaryProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: summaryAddr})

	s3Addr := os.Getenv("S3ADDR")
	if len(s3Addr) == 0 {
		s3Addr = "http://micro-s3:8080"
	}
	s3Proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: s3Addr})

	messageCount := 0
	messageDirector := func(r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			auth = r.URL.Query().Get("auth")
		}
		authUserSessID := sessions.SessionID(strings.TrimPrefix(auth, "Bearer "))
		sessState := &handlers.SessionState{}
		err := context.SessionStore.Get(authUserSessID, sessState)
		if err == nil {
			r.Header.Set("X-User", fmt.Sprintf("{\"userID\":%d}", sessState.User.ID))
		} else {
			r.Header.Del("X-User")
		}

		r.Host = MESSAGESADDR[messageCount]
		r.URL.Host = MESSAGESADDR[messageCount]
		r.URL.Scheme = "http"
		if len(MESSAGESADDR) > messageCount+1 {
			messageCount++
		}
	}

	photoCount := 0
	photoDirector := func(r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			auth = r.URL.Query().Get("auth")
		}
		authUserSessID := sessions.SessionID(strings.TrimPrefix(auth, "Bearer "))
		sessState := &handlers.SessionState{}
		err := context.SessionStore.Get(authUserSessID, sessState)
		if err == nil {
			r.Header.Set("X-User", fmt.Sprintf("{\"userID\":%d}", sessState.User.ID))
		} else {
			r.Header.Del("X-User")
		}

		r.Host = PHOTOSADDR[photoCount]
		r.URL.Host = PHOTOSADDR[photoCount]
		r.URL.Scheme = "http"
		if len(PHOTOSADDR) > photoCount+1 {
			photoCount++
		}
	}

	messageProxy := &httputil.ReverseProxy{Director: messageDirector}
	photoProxy := &httputil.ReverseProxy{Director: photoDirector}

	mux := http.NewServeMux()
	mux.Handle("/v1/summary", summaryProxy)
	mux.Handle("/v1/photos/", photoProxy)
	mux.Handle("/v1/photos", photoProxy)
	mux.Handle("/v1/upload/", s3Proxy)
	// mux.Handle("/v1/channels", messageProxy)
	// mux.Handle("/v1/channels/", messageProxy)
	// mux.Handle("/v1/messages/", messageProxy)
	//messagingProxy := &httputil.ReverseProxy{Director: customDirector(messageURL, context)}
	mux.Handle("/v1/channels", messageProxy)
	mux.Handle("/v1/channels/", messageProxy)
	mux.Handle("/v1/tags", photoProxy)
	mux.Handle("/v1/tags/", photoProxy)
	//mux.Handle("/v1/channels/{channelID}/members", messageProxy)
	mux.Handle("/v1/messages/", messageProxy)
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)
	cors := &handlers.CORS{
		Handle: mux,
	}
	log.Printf("Listening in on Port %s.", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, cors))
}
