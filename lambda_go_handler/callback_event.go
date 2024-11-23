package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
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

	webhookURL = os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		fmt.Println("WEBHOOK_URL is not set")
	}

	data := map[string]string{
		"text": "You said: " + appMentionEvent.Text,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("failed to marshal data")
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("failed to post message to slack")
		return events.LambdaFunctionURLResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to read response body")
		return events.LambdaFunctionURLResponse{}, err
	}

	fmt.Println("status:", resp.StatusCode)
	fmt.Println("responseBody:", string(body))
	return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
}
