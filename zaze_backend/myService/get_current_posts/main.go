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
	claims := request.RequestContext.Authorizer["claims"]
	claimMap := claims.(map[string]interface{})
	UserName := claimMap["cognito:username"].(string)

	svc := s3.New(session.New(), &aws.Config{Region: aws.String("eu-west-2")})

	params := &s3.ListObjectsInput{
		Bucket: aws.String("zaze.io"),
		Prefix: aws.String("user_uploads/static_sites/" + UserName),
	}

	listedObjects, _ := svc.ListObjects(params)
	listedObjectsLength := len(listedObjects.Contents)
	//var strArray [listedObjectsLength]string
	strArray := make([]string, listedObjectsLength)

	for _, key := range listedObjects.Contents {
		// make sure we remove the parts of the s3 bucket the client doesn't need.
		keyItem := strings.Replace(*key.Key, "user_uploads/static_sites/"+UserName+"/", "", -1)
		strArray = append(strArray, keyItem)
	}
	respBody, err := json.Marshal(strArray)
	if err != nil {
		exitErrorf("Could not serialise JSON array.")
	}
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(respBody),
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
