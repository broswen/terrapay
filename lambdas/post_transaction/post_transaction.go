package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {

	// get info from source
	// validate destination exists, get info

	// validate source has amount

	// submit write transaction with 4 records
	// decrement source by amount, validate amount is available
	// increment destination by amount
	// put withdrawal record in source
	// put deposit record in destination

	// return transaction id

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
