package handlers

import (
	"log"
	"net/http"
)

type CORSHandler struct {
	handler http.Handler
}

func (c *CORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("getting cors.")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "600")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		c.handler.ServeHTTP(w, r)
	}
	log.Println("success in cors.")

}

func NewCORS(h http.Handler) *CORSHandler {
	return &CORSHandler{h}
}
