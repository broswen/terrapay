package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {

	// record is INSERT, then it's a new transaction and send notification
	// if type is deposit, send deposit notification to destination
	// if type is withdrawal, send withdrawal notification to source

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
