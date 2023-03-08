package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// Call the ChatGPT API with passed string
func callChatGPT(query string, prompt string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),

		// https://platform.openai.com/docs/guides/chat/chat-vs-completions

		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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
