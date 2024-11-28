package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	modelID = "anthropic.claude-instant-v1"
	region  = "ap-northeast-1"
)

var api *slack.Client
var client *bedrockruntime.Client

func init() {
	token := os.Getenv("BOT_USER_OAUTH_TOKEN")
	api = slack.New(token)

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		fmt.Printf("Error loading AWS configuration: %v\n", err)
		return
	}
	client = bedrockruntime.NewFromConfig(cfg)
}

type Headers struct {
	RetryNum    string `json:"x-slack-retry-num"`
	RetryReason string `json:"x-slack-retry-reason"`
}

type RawRequest struct {
	Body    string  `json:"body"`
	Headers Headers `json:"headers"`
}

func handleRequest(ctx context.Context, lambdaEvent json.RawMessage) (events.LambdaFunctionURLResponse, error) {
	fmt.Printf("Raw lambdaEvent: %v\n", string(lambdaEvent)) // for debug
	var rawRequest *RawRequest
	err := json.Unmarshal(lambdaEvent, &rawRequest)
	if err != nil {
		fmt.Println("failed to unmarshal event")
		return events.LambdaFunctionURLResponse{}, errors.New("failed to unmarshal event")
	}

	/**
	 * This is temporal workaround to avoid duplicated message handling.
	 * map can cause data race and panic, so it should be replaced with a proper solution
	 */
	if (rawRequest.Headers.RetryNum == "1" || rawRequest.Headers.RetryNum == "2") && rawRequest.Headers.RetryReason == "http_timeout" {
		fmt.Println("retry request is ignored")
		return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
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
