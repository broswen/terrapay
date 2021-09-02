package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.KinesisFirehoseEvent) (events.KinesisFirehoseResponse, error) {
	records := make([]events.KinesisFirehoseResponseRecord, 0)
	for _, record := range event.Records {
		responseRecord := events.KinesisFirehoseResponseRecord{
			RecordID: record.RecordID,
			Result:   events.KinesisFirehoseTransformedStateOk,
			Data:     append(record.Data, []byte("\n")...),
		}
		records = append(records, responseRecord)
	}

	return events.KinesisFirehoseResponse{
		Records: records,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
