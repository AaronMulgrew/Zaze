package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	PageName string `json:"PageName"`
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
		PageName: "",
	}

	claims := request.RequestContext.Authorizer["claims"]
	claimMap, claimOK := claims.(map[string]interface{})
	if claimOK == false {
		panic(errors.New("Invalid Credentials"))
	}
	UserName := claimMap["cognito:username"].(string)
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		exitErrorf("Could not decode JSON object. Error: %v", err)
	}
	bucket := aws.String("zaze.io")
	object := aws.String("user_uploads/static_sites/" + UserName + "/" + string(bodyRequest.PageName))
	// BodyRequest will be used to take the json response from client and build it
	log.Print("user_uploads/static_sites/" + UserName + "/" + string(bodyRequest.PageName))
	svc := s3.New(session.New(), &aws.Config{Region: aws.String("eu-west-2")})

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(*bucket), Key: aws.String(*object)})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", object, bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*object),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for object %q to be deleted, %v", *object, err)
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "Deleted OK",
		Headers: map[string]string{
			"Content-Type":                "application/json",
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
