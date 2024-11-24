package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var api *slack.Client

func init() {
	token := os.Getenv("BOT_USER_OAUTH_TOKEN")
	api = slack.New(token)
}

type RawRequest struct {
	Body string `json:"body"`
}

func handleRequest(ctx context.Context, lambdaEvent json.RawMessage) (events.LambdaFunctionURLResponse, error) {
	fmt.Printf("Raw lambdaEvent: %v\n", string(lambdaEvent)) // for debug
	var rawRequest *RawRequest
	err := json.Unmarshal(lambdaEvent, &rawRequest)
	if err != nil {
		fmt.Println("failed to unmarshal event")
		return events.LambdaFunctionURLResponse{}, errors.New("failed to unmarshal event")
	}

	var slackEvent = []byte(rawRequest.Body)

	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(slackEvent),
		slackevents.OptionNoVerifyToken(), // for debug
	)
	if err != nil {
		fmt.Println("failed to parse event")
		return events.LambdaFunctionURLResponse{}, errors.New("failed to parse event")
	}

	fmt.Printf("Parsed Event: %+v\n", eventsAPIEvent) // for debug

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		return handleURLVerification(slackEvent)
	case slackevents.CallbackEvent:
		return handleCallbackEvent(json.RawMessage(slackEvent))
	default:
		fmt.Println("no handler is implemented for this now.")
	}

	return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
