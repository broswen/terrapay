package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.DynamoDBEvent) error {
	for _, record := range event.Records {
		fmt.Println(record)

		if record.EventName != string(events.DynamoDBOperationTypeInsert) {
			continue
		}
		// INSERT types means a new record (either a new account or new transaction)
		newImage := record.Change.NewImage
		switch newImage["type"].String() {
		case "WITHDRAWAL":
			amount, err := newImage["amount"].Float()
			if err != nil {
				fmt.Printf("get amount: %v", err)
				continue
			}
			fmt.Printf("withdrawal from %s to %s for %f\n", newImage["source"].String(), newImage["destination"].String(), amount)
			// send withdrawal notice
		case "DEPOSIT":
			amount, err := newImage["amount"].Float()
			if err != nil {
				fmt.Printf("get amount: %v", err)
				continue
			}
			fmt.Printf("deposit from %s to %s for %f\n", newImage["source"].String(), newImage["destination"].String(), amount)
			// send deposit notice
		default:
			fmt.Println(newImage)
		}
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
