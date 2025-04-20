package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type aiChat interface {
	summarize()
	createClient()
}

type ChatGPT struct {
	ApiKey string
}

type EventsSummary struct {
	Events []struct {
		Title   string `json:"title"`
		Details string `json:"details"`
		Date    string `json:"date"`
	} `json:"events"`
}

func (chat *ChatGPT) SummarizeEvents(events []map[string]string) (EventsSummary, error) {
	client := openai.NewClient(chat.ApiKey)

	var eventsSummary EventsSummary
	schema, err := jsonschema.GenerateSchemaForType(eventsSummary)
	if err != nil {
		log.Fatalf("GenerateSchemaForType error: %v", err)
	}

	jsonBytes, err := json.Marshal(events)
	if err != nil {
		log.Fatal("error encountered while marshalling.")
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a event summarizer that focuses on events happening in Berlin for a given month. Your task is to summarize the given events and details including what is happening (title), details (details) and the start date in mm-dd format (date).",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: string(jsonBytes),
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "events",
					Schema: schema,
					Strict: true,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return eventsSummary, err
	}

	fmt.Printf("Type of content returned: %s", reflect.TypeOf(resp.Choices[0].Message.Content))
	err = schema.Unmarshal(resp.Choices[0].Message.Content, &eventsSummary)
	if err != nil {
		log.Fatalf("Unmarshal schema error: %v", err)
	}

	sort.Slice(eventsSummary.Events, func(i, j int) bool {
		return eventsSummary.Events[i].Date < eventsSummary.Events[j].Date
	})

	return eventsSummary, nil
}
