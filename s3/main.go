package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	AWS_S3_REGION = "us-east-1"
	AWS_S3_BUCKET = "photo-collab"
)

var sess = connectAWS()

func connectAWS() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		panic(err)
	}
	return sess
}

func main() {

	addr, addrExists := os.LookupEnv("ADDR")
	if !addrExists {
		addr = ":8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/upload/", UploadHandler) // Upload
	mux.HandleFunc("/delete/", deleteHandler) // Delete
	log.Printf("Server is open and listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func showError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, message)
}
