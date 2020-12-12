package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//UplaodHandler for uploading img to aws s3 bucket
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	// Get a file from the form input name "file"
	file, header, err := r.FormFile("file")
	if err != nil {
		showError(w, r, http.StatusInternalServerError, "Something went wrong retrieving the file from the form, ")
		fmt.Print(err)
		return
	}
	defer file.Close()

	filename := header.Filename

	// Upload the file to S3.
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET), // Bucket
		Key:    aws.String(filename),      // Name of the file to be saved
		Body:   file,                      // File
	})
	if err != nil {
		//error handling here
		showError(w, r, http.StatusInternalServerError, "Something went wrong uploading the file, ")
		fmt.Print(err)
		return
	}

	fmt.Fprintf(w, "Successfully uploaded to %q\n", AWS_S3_BUCKET)
	return
}

//deleteHandler for deleting imgs from s3 bucket
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	// Get a file from the form input name "file"
	file, header, err := r.FormFile("file")
	if err != nil {
		showError(w, r, http.StatusInternalServerError, "Something went wrong retrieving the file from the form")
		return
	}
	defer file.Close()

	filename := header.Filename

	// create service client.
	delete := s3.New(sess)

	_, err = delete.DeleteObject(&s3.DeleteObjectInput{
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

	fmt.Fprintf(w, "Successfully deleted image %q from bucket to %q\n", filename, AWS_S3_BUCKET)
	return
}
