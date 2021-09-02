package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	event := events.APIGatewayV2HTTPRequest{
		Body: `{"password":"password","email":"test@test.com"}`,
	}
	response, err := HandleRequest(context.TODO(), event)
	fmt.Println(response.Body)
	if err != nil {
		t.Errorf("HandleRequest failed: %w", err)
	}
}
