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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	guuid "github.com/google/uuid"
)

// Item Create struct to hold info about new item
type Item struct {
	UniqueID  string
	PostTitle string
	UserName  string
}

// AddToDynamoDB is a function to export post details to DynamoDB
func AddToDynamoDB(UniqueID string, Title string, UserName string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)
	svc := dynamodb.New(sess)

	if err != nil {
		exitErrorf("Could not create AWS Session.")
	}

	item := Item{
		UniqueID:  UniqueID,
		PostTitle: Title,
		UserName:  UserName,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tableName := "zaze-user-posts"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
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
	claims := request.RequestContext.Authorizer["claims"]
	claimMap := claims.(map[string]interface{})
	UserName := claimMap["cognito:username"].(string)

	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		HTMLContents: "",
		Title:        "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		exitErrorf("Could not decode JSON object. Error: %v", err)
	}
	uniqueID := guuid.New().String()

	bucket := "zaze.io"
	filename := "user_uploads/static_sites/" + UserName + "/" + uniqueID + ".html"

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
	AddToDynamoDB(uniqueID, bodyRequest.Title, UserName)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "https://users.zaze.io/" + UserName + "/" + uniqueID + ".html",
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
