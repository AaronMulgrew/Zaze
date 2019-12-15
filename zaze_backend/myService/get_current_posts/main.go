package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Item Create struct to hold info about new item
type Item struct {
	UniqueID  string
	PostTitle string
	UserName  string
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	// BodyRequest will be used to take the json response from client and build it
	claims := request.RequestContext.Authorizer["claims"]
	claimMap, claimOK := claims.(map[string]interface{})
	if claimOK == false {
		panic(errors.New("Invalid Credentials"))
	}
	UserName := claimMap["cognito:username"].(string)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)
	if err != nil {
		exitErrorf("Unable to connect to aws with error %v", err)
	}
	svc := dynamodb.New(sess)
	tableName := "zaze-user-posts"

	filt := expression.Name("UserName").Equal(expression.Value(UserName))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// snippet-end:[dynamodb.go.scan_items.expr]

	// snippet-start:[dynamodb.go.scan_items.call]
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	var allItems []Item

	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		allItems = append(allItems, item)
	}

	respBody, errJSON := json.Marshal(allItems)
	if errJSON != nil {
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
