package main

import (
	"context"
	"fmt"
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
	RequestHeader          string `json:"header"`
	RequestContent         string `json:"content"`
	RequestBackgroundColor string `json:"background-color"`
	RequestFontColor       string `json:"font-color"`
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	bucket := "zaze.io"
	filename := "user_uploads/static_sites/xxxx/test.txt"

	// this is your data that you have in memory
	// in this example it is hard coded but it may come from very distinct
	// sources, like streaming services for example.
	data := "Hello, world!"

	// create a reader from data data in memory
	reader := strings.NewReader(data)

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
		Body: reader,
	})
	if err != nil {
		exitErrorf("Zero bytes returned error. Error: %v", err)
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "OK",
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
