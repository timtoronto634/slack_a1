package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack/slackevents"
)

// handleURLVerification handles the URL verification event
// it is used when setting up request URL in Event Subscriptions
func handleURLVerification(slackEvent []byte) (events.LambdaFunctionURLResponse, error) {
	var r *slackevents.EventsAPIURLVerificationEvent
	err := json.Unmarshal(slackEvent, &r)
	if err != nil {
		fmt.Println("failed to unmarshal url verification event")
		return events.LambdaFunctionURLResponse{}, errors.New("failed to unmarshal url verification event")
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       r.Challenge,
	}, nil
}
