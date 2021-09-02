package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	var request RegisterRequest
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	fmt.Printf("%v\n", request)

	// validate request

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	fmt.Println(string(hashedBytes))
	// check if user already exists
	// bcrypt hash password
	// insert into dynamodb table

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
