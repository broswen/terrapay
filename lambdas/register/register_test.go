package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	event := events.APIGatewayV2HTTPRequest{
		Body: `{"password":"password","email":"test@test.com"}`,
	}
	_, err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("HandleRequest failed: %w", err)
	}
}
