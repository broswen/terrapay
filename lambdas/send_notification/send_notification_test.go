package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	event := events.DynamoDBEvent{
		Records: []events.DynamoDBEventRecord{
			{
				EventName: "INSERT",
				Change: events.DynamoDBStreamRecord{
					NewImage: map[string]events.DynamoDBAttributeValue{
						"type": events.DynamoDBAttributeValue{
							"DEPOSIT",
							events.DataTypeString,
						},
						"amount": events.DynamoDBAttributeValue{
							"123.45",
							events.DataTypeNumber,
						},
						"source": events.DynamoDBAttributeValue{
							"AccountA",
							events.DataTypeString,
						},
						"destination": events.DynamoDBAttributeValue{
							"AccountB",
							events.DataTypeString,
						},
					},
				},
			},
		},
	}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("HandleRequest failed: %w", err)
	}
}
