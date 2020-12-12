package handlers

import (
	"assignments-ydf1014/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket

type SocketStore struct {
	Connections map[int64]*websocket.Conn
	// Sh          *Context
	// Connections []*websocket.Conn
	lock sync.Mutex
}

// Control messages for websocket

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// Thread-safe method for inserting a connection
func (ws *SocketStore) InsertConnection(conn *websocket.Conn, userID int64) {
	log.Println("In InsertConnection")
	ws.lock.Lock()
	// connID := len(ws.Connections)
	// insert socket connection
	ws.Connections[userID] = conn

	ws.lock.Unlock()
}

// Thread-safe method for inserting a connection
func (ws *SocketStore) RemoveConnection(userID int64) {
	ws.lock.Lock()
	// insert socket connection
	delete(ws.Connections, userID)
	// s.Connections = append(s.Connections[:connId], s.Connections[connId+1:]...)
	ws.lock.Unlock()
}

// Simple method for writing a message to all live connections.
// In your homework, you will be writing a message to a subset of connections
// (if the message is intended for a private channel), or to all of them (if the message
// is posted on a public channel
// func (ws *SocketStore) WriteToAllConnections(message []byte) error {
// 	var writeError error

// 	for _, conn := range ws.Connections {
// 		// messageType is 1
// 		// data was the message
// 		writeError = conn.WriteMessage(1, message)
// 		if writeError != nil {
// 			return writeError
// 		}
// 	}

// 	return nil
// }

type rabbitMessage struct {
	UserIDs []int64 `json:"userIDs, omitempty"`
}

func (ws *SocketStore) SendMessages(messages <-chan amqp.Delivery) {
	for msg := range messages {
		// msg := <-messages
		// log.Println("In SendMessages")
		// ws.lock.Lock()
		// log.Println("In SendMessages 2")
		userIDs := &rabbitMessage{}
		err := json.Unmarshal(msg.Body, userIDs)
		if err != nil {
			log.Printf("Error unmarshalling userIDs %v", err)
		}
		// log.Println(userIDs)
		if len(userIDs.UserIDs) == 0 {
			for userID, conn := range ws.Connections {
				log.Printf("Message 1: %v", msg.Body)
				err := conn.WriteMessage(TextMessage, msg.Body)
				if err != nil {
					log.Printf("fail to write message: %s", err)
					ws.RemoveConnection(userID)
				}
			}
		} else {
			for _, userID := range userIDs.UserIDs {
				conn := ws.Connections[userID]
				// conn, errBool := ws.Connections[userID]
				if conn != nil {
					err = conn.WriteMessage(TextMessage, msg.Body)
					if err != nil {
						log.Printf("cannot writemessage in else branch: %s", err)
						ws.RemoveConnection(userID)
						conn.Close()
					}
				} else {
					log.Printf("conn is null")
				}
				// if !errBool {
				// 	ws.RemoveConnection(userID)
				// }
				// log.Printf("Message 2: %v", msg.Body)
				// err := conn.WriteMessage(TextMessage, msg.Body)
				// if err != nil {
				// 	ws.RemoveConnection(userID)
				// 	conn.Close()
				// }
			}
		}
		// log.Println("In SendMessages 3")
		// ws.lock.Unlock()
	}
	// log.Println("In SendMessages 3")
	// ws.lock.Unlock()
	// return nil
}

// // This is a struct to read our message into
// type msg struct {
// 	Message string `json:"message"`
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// if r.Header.Get("Origin") == "https://nehay.me" {
		return true
		// }
		// This function's purpose is to reject websocket upgrade requests if the
		// origin of the websockete handshake request is coming from unknown domains.
		// This prevents some random domain from opening up a socket with your server.
		// TODO: make sure you modify this for your HW to check if r.Origin is your host
		// fmt.Sprintf("Connection Refused", 403)
		// return false
	},
}

func NewSocketStorego() *SocketStore {
	return &SocketStore{
		// Sh:          context,
		Connections: make(map[int64]*websocket.Conn),
	}
}

func (context *Context) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// handle the websocket handshake
	// log.Println(ws.Sh)

	// add session state
	state := &SessionState{}
	// log.Println(ws.Sh.SigningKey)
	// log.Println(ws.Sh.Store)
	_, err := sessions.GetState(r, context.Key, context.Session, state)
	// log.Println(state.User)
	// log.Println(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("UserHandler: error getting session state/session unauthorized %v", err),
			http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}
	// Insert our connection onto our datastructure for ongoing usage
	//   connId := ws.InsertConnection(conn, state.User.ID)
	context.SocketStore.InsertConnection(conn, state.User.ID)

	// Invoke a goroutine for handling control messages from this connection
	go (func(conn *websocket.Conn, userID int64) {
		defer conn.Close()
		defer context.SocketStore.RemoveConnection(userID)

		for {
			messageType, p, err := conn.ReadMessage()

			if messageType == TextMessage || messageType == BinaryMessage {
				fmt.Printf("Client says %v\n", p)
				// fmt.Printf("Writing %s to all sockets\n", string(p))
				// ws.WriteToAllConnections(append([]byte("Hello from server: "), p...))
			} else if messageType == CloseMessage {
				fmt.Println("Close message received.")
				break
			} else if err != nil {
				fmt.Println("Error reading message.")
				break
			}
			// ignore ping and pong messages
		}

	})(conn, state.User.ID)
}
