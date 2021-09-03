package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {

	// get info from account

	// parse query params ?from=2021-08-12&to=2021-08-18

	// return account balance, id
	// return list of transactions within date range

	response := struct {
		Msg string `json:"msg"`
	}{
		Msg: "hello",
	}

	j, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func main() {
	lambda.Start(HandleRequest)
}
