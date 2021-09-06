package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/broswen/terrapay/account"
)

func HandleRequest(ctx context.Context, event events.DynamoDBEvent) error {
	for _, record := range event.Records {
		fmt.Println(record)

		if record.EventName != string(events.DynamoDBOperationTypeInsert) {
			continue
		}
		// INSERT types means a new record (either a new account or new transaction)
		newImage := record.Change.NewImage
		switch newImage["action"].String() {
		case string(account.Pay):
			amount, err := newImage["amount"].Float()
			id := newImage["id"]
			sender := newImage["sender"]
			recipient := newImage["recipient"]
			timestamp := newImage["timestamp"]
			if err != nil {
				fmt.Printf("get amount: %v", err)
				continue
			}
			fmt.Printf("Sending payment from %s to %s\nId: %s\nAmount: %.2f\nTimestamp: %s\n", sender, recipient, id, amount, timestamp)
		case string(account.Receive):
			amount, err := newImage["amount"].Float()
			id := newImage["id"]
			sender := newImage["sender"]
			recipient := newImage["recipient"]
			timestamp := newImage["timestamp"]
			if err != nil {
				fmt.Printf("get amount: %v", err)
				continue
			}
			fmt.Printf("Sending payment from %s to %s\nId: %s\nAmount: %.2f\nTimestamp: %s\n", sender, recipient, id, amount, timestamp)
		default:
			fmt.Println(newImage)
		}
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
