package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

	// check if user account exists

	// if exists, hash and compare password hashes

	// if valid user, generate jwt and return

	response := LoginResponse{
		Token: "test",
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

func main() {
	lambda.Start(HandleRequest)
}
