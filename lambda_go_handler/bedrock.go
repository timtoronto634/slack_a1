package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/samber/lo"
)

type SingleMessage struct {
	Sender string
	Text   string
}

type AnthropicRequest struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature"`
	TopP              float64  `json:"top_p"`
	TopK              int      `json:"top_k"`
	StopSequences     []string `json:"stop_sequences"`
}

type AnthropicResponse struct {
	Completion string `json:"completion"`
}

func summarizeConversation(messages []SingleMessage) string {
	conversation := lo.Map(messages, func(m SingleMessage, _ int) string {
		return fmt.Sprintf("%s: %s\n", m.Sender, m.Text)
	})

	prompt := fmt.Sprintf(`
Human:
Summarize the following Conversation while following the restrictions 
### restrictions ### 
- The summarization must be within 200 words. 
- The summarization must be in the same language as the conversation
- The summarization must be completed.
- The summarization immediately start without leading text 

### conversion ###
%s

Assistant:
`, conversation)

	completion, err := callBedrock(prompt)
	if err != nil {
		fmt.Printf("error calling bedrock: %v\n", err)
		return ""
	}

	return completion
}

func callBedrock(prompt string) (string, error) {
	request := AnthropicRequest{
		Prompt:            prompt,
		MaxTokensToSample: 4000,
		Temperature:       0.7,
		TopP:              0.9,
		TopK:              50,
		StopSequences:     []string{},
	}

	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	output, err := client.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		Body:        jsonPayload,
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return "", fmt.Errorf("error invoking model: %w", err)
	}

	var response AnthropicResponse
	err = json.Unmarshal(output.Body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response.Completion, nil
}
