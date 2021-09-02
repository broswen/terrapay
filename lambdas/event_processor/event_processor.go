package main

import (
	"context"

	"encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.KinesisFirehoseEvent) (events.KinesisFirehoseResponse, error) {
	records := make([]events.KinesisFirehoseResponseRecord, 0)
	for _, record := range event.Records {

		var decoded []byte
		_, err := base64.StdEncoding.Decode(decoded, record.Data)
		if err != nil {
			return events.KinesisFirehoseResponse{}, err
		}
		var encoded []byte
		base64.StdEncoding.Encode(encoded, append(decoded, []byte("\n")...))

		responseRecord := events.KinesisFirehoseResponseRecord{
			RecordID: record.RecordID,
			Result:   events.KinesisFirehoseTransformedStateOk,
			Data:     encoded,
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
