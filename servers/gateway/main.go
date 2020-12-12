package main

import (
	"assignments-ydf1014/servers/gateway/handlers"
	"assignments-ydf1014/servers/gateway/models/users"
	"assignments-ydf1014/servers/gateway/sessions"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

//main is the main entry point for the server
func main() {

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	// messageAddrs := strings.Split(os.Getenv("MESSAGE_ADDR"), ", ")
	messageURLArr := strings.Split(os.Getenv("MESSAGE_ADDR"), ", ")

	// var messageURLArr []*url.URL
	// for _, address := range messageAddrs {
	// 	messageURL, err := url.Parse(address)
	// 	if err != nil {
	// 		fmt.Sprintf("cannot parse message addresss: %v", err)
	// 		return
	// 	}
	// 	messageURLArr = append(messageURLArr, messageURL)
	// }

	summaryURLArr := strings.Split(os.Getenv("SUMMARY_ADDR"), ", ")

	sessionKey := os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		fmt.Println("please set SESSIONKEY")
		os.Exit(1)
	}

	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		fmt.Println("please set REDISADDR")
		os.Exit(1)
	}

	tlsKey := os.Getenv("TLSKEY")
	if len(tlsKey) == 0 {
		fmt.Println("please set TLS KEY")
		os.Exit(1)
	}
	tlsCert := os.Getenv("TLSCERT")
	if len(tlsCert) == 0 {
		fmt.Println("please set TLS CERT")
		os.Exit(2)
	}

	redis := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	redisStore := sessions.NewRedisStore(redis, time.Hour)

	dsn := os.Getenv("DSN")
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v", err)
		os.Exit(1)
	}
	user := users.NewMySQLStorego(database)
	socketStore := handlers.NewSocketStorego()
	context := handlers.Initialize(redisStore, user, sessionKey, socketStore)

	messageProxy := NewServiceProxy(messageURLArr, sessionKey, redisStore)
	summaryProxy := NewServiceProxy(summaryURLArr, sessionKey, redisStore)
	mux := http.NewServeMux()
	//WEBSOCKET /////////////////////////////////////////////
	///////////////////////////////////////
	/////////////////////////////////////////////
	///////////////////////////////////////////////////
	rabbitAddr := os.Getenv("RABBIT_ADDR")
	if rabbitAddr == "" {
		fmt.Printf("please provide rabbit addr")
		os.Exit(1)
	}
	conn, err := amqp.Dial("amqp://mq:5672/")
	if err != nil {
		fmt.Printf("fail connecting to rabbitmq: %s", err)
		os.Exit(1)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("fail opening rabbit channel: %s", err)
		os.Exit(1)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitAddr, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		fmt.Printf("fail queue declare: %s", err)
		os.Exit(1)
	}
	msg, err := ch.Consume(
		q.Name,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
	if err != nil {
		fmt.Printf("fail consume: %s", err)
		os.Exit(1)
	}
	go context.SocketStore.SendMessages(msg)
	//WEBSOCKET /////////////////////////////////////////////
	///////////////////////////////////////
	/////////////////////////////////////////////
	///////////////////////////////////////////////////
	mux.HandleFunc("/v1/ws", context.WebSocketConnectionHandler)
	mux.Handle("/v1/summary", summaryProxy)
	mux.Handle("/v1/messages", messageProxy)
	mux.Handle("/v1/messages/", messageProxy)
	mux.Handle("/v1/channels", messageProxy)
	mux.Handle("/v1/channels/", messageProxy)
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)
	// mux.Handle("/v1/ws", handlers.NewWebSocketHandler(ctx))
	middlewareWrapCORS := handlers.NewCORS(mux)
	log.Printf("server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, middlewareWrapCORS))

}

func NewServiceProxy(targets []string, signingKey string, store sessions.Store) *httputil.ReverseProxy {
	counter := 0
	state := &handlers.SessionState{}
	mx := sync.Mutex{}
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			mx.Lock()

			targ := targets[counter]
			targets32 := len(targets)

			counter = (counter + 1) % targets32

			mx.Unlock()
			log.Printf("reverse proxy get state")

			r.Header.Del("X-User")

			log.Printf("reverse proxy get marshal")
			j, err := json.Marshal(state.User)
			if err != nil {
				fmt.Sprintf("cannot get marshal %v", err)
				return
			}
			r.Header.Add("X-User", string(j))
			r.Host = targ
			r.URL.Host = targ
			log.Printf("host in request url is:")
			log.Printf(r.URL.Host)
			r.URL.Scheme = "http"
			log.Printf("done return reverse proxy")
		},
	}
}
