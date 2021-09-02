package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	event := events.KinesisFirehoseEvent{
		InvocationID:           "123",
		DeliveryStreamArn:      "abc",
		SourceKinesisStreamArn: "def",
		Records: []events.KinesisFirehoseEventRecord{
			{
				RecordID: "1",
				Data:     []byte("SGVsbG8sIHRoaXMgaXMgYSB0ZXN0IDEyMy4="),
			},
		},
	}
	_, err := HandleRequest(event)
	if err != nil {
		t.Errorf("HandleRequest failed: %w", err)
	}
}
