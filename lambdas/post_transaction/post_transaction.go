package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/broswen/terrapay/account"
)

var ddbClient *dynamodb.Client
var accountService *account.AccountService

type PostTransactionRequest struct {
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

type PostTransactionResponse struct {
	Transactions account.Transaction
}

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	authorization, ok := event.Headers["authorization"]
	if !ok {
		log.Println("missing jwt")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 401,
		}, nil
	}

	authorizationParts := strings.Split(authorization, " ")
	if len(authorizationParts) != 2 {
		log.Println("invalid parts for bearer authorization")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 401,
		}, nil
	}

	claims, err := account.ValidateJWT(authorizationParts[1])
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 401,
		}, nil
	}

	var request PostTransactionRequest
	err = json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	fmt.Printf("%v\n", request)

	transaction, err := accountService.PostTransaction(ctx, claims.Subject, request.Recipient, request.Amount)
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
		}, nil
	}
	// validate source has amount

	// submit write transaction with 4 records
	// decrement source by amount, validate amount is available
	// increment destination by amount
	// put withdrawal record in source
	// put deposit record in destination

	// return transaction id

	response := PostTransactionResponse{
		Transactions: transaction,
	}

	j, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(j),
	}, nil
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	ddbClient = dynamodb.NewFromConfig(cfg)
	accountService, err = account.NewFromClient(ddbClient)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
