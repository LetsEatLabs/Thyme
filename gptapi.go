package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

/////////////////

// Call the ChatGPT API with passed string and using a prompt
func callChatGPT(query string, prompt string, model string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),

		// https://platform.openai.com/docs/guides/chat/chat-vs-completions

		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	return resp.Choices[0].Message.Content
}

/////////////////

// Call the ChatGPT API with passed string and no prompt
func callChatGPTNoPrompt(query string, model string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),

		// https://platform.openai.com/docs/guides/chat/chat-vs-completions

		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	return resp.Choices[0].Message.Content
}

/////////////////

// Call the GPT Completions API
func callGPT(query string) string {
	c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Davinci,
		MaxTokens: 250,
		Prompt:    query,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	return resp.Choices[0].Text
}

/////////////////

// Save the GPT Completions API response to a file
// If the THYME_QUERY_LOGGING_DIR environment variable is not set, do nothing
// Timestamp is when it saves, not when you send the query.
func saveGPT(qs QuerySave) {
	saveDir := os.Getenv("THYME_QUERY_LOGGING_DIR")

	if saveDir == "" {
		return
	}

	// Write the file as a json sting {'prompt': '...', 'query': '...', 'answer': ...}
	// Filename is YYYY-MM-DD-HH-mm-SS-query.json
	currentTime := time.Now()
	year, month, day := currentTime.Date()
	hour, min, sec := currentTime.Clock()
	formattedTime := fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d", year, month, day, hour, min, sec)
	formattingTimeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)

	filename := fmt.Sprintf("%s/%s-query.json", saveDir, formattedTime)

	fileData := fmt.Sprintf("{'timestamp': '%s', 'prompt': '%s', 'query': '%s', 'answer': '%s'}", formattingTimeStamp, qs.Prompt, qs.Query, qs.Answer)

	// Write the file
	err := ioutil.WriteFile(filename, []byte(fileData), 0644)

	if err != nil {
		fmt.Println(err)
	}
}
