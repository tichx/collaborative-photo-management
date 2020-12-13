package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	AWS_S3_REGION = "us-west-2"
	AWS_S3_BUCKET = "image-441"
)

var sess = connectAWS()

func connectAWS() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION), Credentials: credentials.NewStaticCredentials("AKIAJOMUO3S2R36HG3JQ", "tdfNki33dgN9EcDki2tMjp0ToW2SE6BJZvK4omoV", "")})
	if err != nil {
		panic(err)
	}
	return sess
}

func main() {

	addr, addrExists := os.LookupEnv("ADDR")
	if !addrExists {
		addr = ":8181"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/upload/", UploadHandler) // Upload
	mux.HandleFunc("/v1/delete/", deleteHandler) // Delete
	log.Printf("Server is open and listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func showError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, message)
}
