package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	UserName     string `json:"username"`
	HTMLContents string `json:"HTMLContents"`
	Title        string `json:"title"`
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		UserName:     "",
		HTMLContents: "",
		Title:        "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		exitErrorf("Could not decode JSON object. Error: %v", err)
	}

	log.Print("HTMLContents: " + bodyRequest.HTMLContents)
	bucket := "zaze.io"
	filename := "user_uploads/static_sites/" + bodyRequest.UserName + "/" + bodyRequest.Title + ".html"

	// create a reader from data data in memory
	reader := strings.NewReader(bodyRequest.HTMLContents)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)
	if err != nil {
		exitErrorf("Could not create AWS Session.")
	}
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		// here you pass your reader
		// the aws sdk will manage all the memory and file reading for you
		Body:        reader,
		ContentType: aws.String("text/html"),
	})
	if err != nil {
		exitErrorf("S3 upload error. Error: %v", err)
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "https://users.zaze.io/" + bodyRequest.UserName + "/" + bodyRequest.Title + ".html",
		Headers: map[string]string{
			"Content-Type":                "text/html",
			"Access-Control-Allow-Origin": "*",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
