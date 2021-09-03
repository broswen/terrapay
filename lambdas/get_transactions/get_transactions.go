package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// get info from account

	// parse query params ?from=2021-08-12&to=2021-08-18

	// return account balance, id
	// return list of transactions within date range

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
