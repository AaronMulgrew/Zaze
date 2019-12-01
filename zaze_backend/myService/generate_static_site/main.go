package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		RequestHeader:          "",
		RequestContent:         "",
		RequestBackgroundColor: "",
		RequestFontColor:       "",
	}

	//var test string
	//test = "<html><h1>hello world</h1></html>"

	//body, err := json.Marshal(map[string]interface{}{
	//	"message": test,
	//})
	buff := &aws.WriteAtBuffer{}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		exitErrorf("Could not decode JSON object.")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)
	if err != nil {
		exitErrorf("Unable to connect to aws with error %v", err)
	}

	downloader := s3manager.NewDownloader(sess)

	// Create S3 service client
	bucket := "zaze-templates"
	item := "hello.html"

	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}
	if numBytes == 0 {
		exitErrorf("Zero bytes returned error.")
	}

	data := buff.Bytes() // now data is my []byte array
	strData := string(data)
	result := strings.Replace(strData, "<headervalue>", bodyRequest.RequestHeader, -1)
	resultWithContent := strings.Replace(result, "<contentvalue>", bodyRequest.RequestContent, -1)
	resultWithCSS := strings.Replace(resultWithContent, "<cssvalue>", "background-color:"+bodyRequest.RequestBackgroundColor, -1)
	resultWithCSSandFont := strings.Replace(resultWithCSS, "<font-color>", "color:"+bodyRequest.RequestFontColor, -1)

	//for _, b := range result.Buckets {
	//		fmt.Printf("* %s created on %s\n",
	//				aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	//}
	body := resultWithCSSandFont
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
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
