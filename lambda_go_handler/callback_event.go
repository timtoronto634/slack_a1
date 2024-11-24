package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func handleCallbackEvent(slackEvent json.RawMessage) (events.LambdaFunctionURLResponse, error) {
	fmt.Printf("Raw Event: %v\n", string(slackEvent)) // for debug

	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(slackEvent),
		slackevents.OptionNoVerifyToken(), // for debug
	)
	if err != nil {
		fmt.Println("failed to parse event")
		return events.LambdaFunctionURLResponse{}, errors.New("failed to parse event")
	}
	fmt.Printf("Parsed eventsAPIEvent: %+v\n", eventsAPIEvent) // for debug

	if eventsAPIEvent.InnerEvent.Type != string(slackevents.AppMention) {
		fmt.Printf("only mention is handled. event type: %s\n", eventsAPIEvent.InnerEvent.Type)
		return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
	}

	var eventsAPICallbackEvent *slackevents.EventsAPICallbackEvent
	err = json.Unmarshal(slackEvent, &eventsAPICallbackEvent)
	if err != nil {
		fmt.Println("failed to unmarshal event")
	}

	appMentionType := reflect.TypeOf(slackevents.AppMentionEvent{})
	recvEvent := reflect.New(appMentionType).Interface()
	err = json.Unmarshal(*eventsAPICallbackEvent.InnerEvent, recvEvent)
	if err != nil {
		fmt.Println("failed to unmarshal event")

		return events.LambdaFunctionURLResponse{}, err
	}

	fmt.Printf("recvEvent: %+v\n", recvEvent)
	appMentionEvent := recvEvent.(*slackevents.AppMentionEvent)

	fmt.Printf("appMentionEvent: %+v\n", appMentionEvent)

	api.PostMessage(appMentionEvent.Channel, slack.MsgOptionText("by app: "+appMentionEvent.Text, false))
	return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
}
