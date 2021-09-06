package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/broswen/terrapay/account"
)

var ddbClient *dynamodb.Client
var accountService *account.AccountService

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	var request LoginRequest
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	fmt.Printf("%v\n", request)

	// validate request

	token, err := accountService.Login(ctx, request.Email, request.Password)

	if err != nil {
		log.Fatal(err)
	}

	response := LoginResponse{
		Token: token,
	}
	j, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
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
