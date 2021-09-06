package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/broswen/terrapay/account"
)

type PostTransactionRequest struct {
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

type PostTransactionResponse struct {
	Transaction account.Transaction
}

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// get info from source
	// validate destination exists, get info

	// validate source has amount

	// submit write transaction with 4 records
	// decrement source by amount, validate amount is available
	// increment destination by amount
	// put withdrawal record in source
	// put deposit record in destination

	// return transaction id

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
