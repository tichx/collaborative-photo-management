package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//ImgStruct has meta-data for image
type ImgStruct struct {
	ImgUrl       string    `json:"imgurl"`
	DateModified time.Time `json:"datemodified"`
}

//UploadHandler for uploading img to aws s3 bucket
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	/*
		- Get the `url` query string parameter value from the request.
		  If not supplied, respond with an http.StatusBadRequest error.
	*/
	//queryString, ok := r.URL.Query()["key"]
	//filename := r.URL.Query().Get("filename")

	uploadContentType := "image/jpeg"
	acl := "public-read"
	contentDisposition := "inline"

	r.ParseMultipartForm(10 << 20)

	// Get a file from the form input name "file"
	file, header, err := r.FormFile("file")
	if err != nil {
		showError(w, r, http.StatusInternalServerError, "Something went wrong retrieving the file from the form, ")
		fmt.Print(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", header.Filename)
	fmt.Printf("File Size: %+v\n", header.Size)
	fmt.Printf("MIME Header: %+v\n", header.Header)

	filename := header.Filename

	// Upload the file to S3.
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(AWS_S3_BUCKET), // Bucket
		Key:                aws.String(filename),      // Name of the file to be saved
		Body:               file,
		ContentDisposition: aws.String(contentDisposition), //
		ContentType:        aws.String(uploadContentType),  // this is what you need!
		ACL:                aws.String(acl),                // this makes it public so people can see it
	})
	if err != nil {
		//error handling here
		showError(w, r, http.StatusInternalServerError, "Something went wrong uploading the file, ")
		fmt.Print(err)
		return
	}

	data := ImgStruct{"https://image-441.s3.amazonaws.com/" + filename, time.Now()}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		showError(w, r, http.StatusInternalServerError, "Error getting img metadata")
		return
	}

	return
}

//deleteHandler for deleting imgs from s3 bucket
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	filename := r.URL.Query().Get("imgname")

	// create service client.
	delete := s3.New(sess)

	_, err := delete.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET), // Bucket
		Key:    aws.String(filename),      // Name of the file to be deleted
	})
	if err != nil {
		//error handling here
		showError(w, r, http.StatusInternalServerError, "Unable to delete object "+filename+" from bucket "+AWS_S3_BUCKET+", ")
		fmt.Print(err)
		return
	}
	err = delete.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(filename),
	})

	if err != nil {
		showError(w, r, http.StatusInternalServerError, "Error occured while waiting for object "+filename+" from bucket "+AWS_S3_BUCKET+" to be deleted ,")
		fmt.Print(err)
		return
	}

	data := ImgStruct{"https://image-441.s3.amazonaws.com/" + filename, time.Now()}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		showError(w, r, http.StatusInternalServerError, "Error getting img metadata")
		return
	}

	return
}
