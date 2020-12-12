package main

import (
	"assignments-ydf1014/servers/summary/handlers"
	"log"
	"net/http"
	"os"
)

//main is the main entry point for the server
func main() {
	/* TODO: add code to do the following
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

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}

	// // tlsKey := os.Getenv("TLSKEY")
	// tlsKey := "/etc/letsencrypt/live/chatapi.danfengy.me/privkey.pem"
	// if len(tlsKey) == 0 {
	// 	fmt.Println("please set TLS KEY")
	// 	os.Exit(1)
	// }

	// //tlsCert := os.Getenv("TLSCERT")

	// tlsCert := "/etc/letsencrypt/live/chatapi.danfengy.me/fullchain.pem"

	// if len(tlsCert) == 0 {
	// 	fmt.Println("please set TLS CERT")
	// 	os.Exit(2)
	// }

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	log.Printf("server is listening at http://%s", addr)
	// log.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, mux))
	log.Fatal(http.ListenAndServe(addr, mux))

}
