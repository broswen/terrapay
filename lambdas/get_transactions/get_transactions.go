package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {
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
