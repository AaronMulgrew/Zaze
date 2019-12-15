package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	PostName string `json:"postname"`
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
	claimMap, claimOK := claims.(map[string]interface{})
	if claimOK == false {
		panic(errors.New("Invalid Credentials"))
	}
	UserName := claimMap["cognito:username"].(string)
	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		PostName: "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		exitErrorf("Could not decode JSON object. Error: %v", err)
	}

	svc := s3.New(session.New(), &aws.Config{Region: aws.String("eu-west-2")})

	input := &s3.GetObjectInput{
		Bucket: aws.String("zaze.io"),
		Key:    aws.String("user_uploads/static_sites/" + UserName + "/" + bodyRequest.PostName + ".html"),
	}
	log.Print("user_uploads/static_sites/" + UserName + "/" + bodyRequest.PostName + ".html")
	result, err := svc.GetObject(input)
	if err != nil {
		panic(errors.New("no object found"))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)
	htmlPost := buf.String()

	start, stop := "{{", "}}" // just replace these with whatever you like...
	sSplits := strings.Split(htmlPost, start)
	embeddedElements := []string{}

	if len(sSplits) > 1 { // n splits = 1 means start char not found!
		for _, subStr := range sSplits { // check each substring for end
			ixEnd := strings.Index(subStr, stop)
			if ixEnd != -1 {
				embeddedElements = append(embeddedElements, subStr[:ixEnd])
			}
		}
	}

	log.Print(embeddedElements)
	embeddedElementsString := strings.Join(embeddedElements, ",") // join the elements
	embeddedElementsString = "{" + embeddedElementsString + "}"
	// but do not add anything to the array
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            embeddedElementsString,
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
