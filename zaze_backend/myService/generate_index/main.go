package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

// AddToS3 is a functon which simply adds that file to S3.
func AddToS3(HTMLContents string, UserName string) {
	bucket := "zaze.io"
	filename := "user_uploads/static_sites/" + UserName + "/" + "index.html"

	// create a reader from data data in memory
	reader := strings.NewReader(HTMLContents)

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
		Body:         reader,
		ContentType:  aws.String("text/html"),
		CacheControl: aws.String("no-cache"),
	})
	if err != nil {
		exitErrorf("S3 upload error. Error: %v", err)
	}
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

	HTMLContents := `<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
		<style>
		body, html {
			background-color: #05111e !important;
			color: #AAA !important;
		}

		a {
			color: #FFF !important;

		}

		a:hover{
			text-decoration: none !important;
		}

		.navbar {
			margin-bottom: 0px;
		}
		.jumbotron{
			min-height: 100%;  /* Fallback for browsers do NOT support vh unit */
			min-height: 100vh; /* These two lines are counted as one :-)       */

			min-height: calc(100% - 64px);  /*  */
			min-height: calc(100vh - 64px); /* */

			display: flex;
			align-items: center;
			background-color: #05111e !important;
		}



	  </style>
    <title>` + UserName + `'s Posts</title>
  </head>
  <body>
	<div class="jumbotron d-flex align-items-center">
		<div class="container text-center">
			<h1>` + UserName + `'s Posts</h1>
	`

	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}

		HTMLContents = HTMLContents + "<h2><a href=\"" + item.UniqueID + ".html\">" + item.PostTitle + "</a></h2>"
		//allItems = append(allItems, item)
	}
	HTMLContents = HTMLContents + `</div></div><script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
  </body>
</html>`
	AddToS3(HTMLContents, UserName)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "OK",
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
