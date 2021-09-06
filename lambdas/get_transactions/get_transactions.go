package main

import (
	"context"
	"encoding/json"
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

type GetTransactionsResponse struct {
	Transactions []account.Transaction `json:"transactions"`
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

	// TODO parse query params ?from=2021-08-12&to=2021-08-18

	transactions, err := accountService.GetTransactions(ctx, claims.Subject)
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
		}, nil
	}

	response := GetTransactionsResponse{
		Transactions: transactions,
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
