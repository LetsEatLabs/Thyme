package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

// Call the ChatGPT API with passed string and using a prompt and passing a JSON schema, no prompt
func callChatGPTFunctionCall(query string, prompt string, model string, funcCall []byte) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	var FunctionCallObj openai.FunctionCall
	FunctionCallObj.Name = "functioncall"
	FunctionCallObj.Arguments = string(funcCall)

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
					Role:         openai.ChatMessageRoleAssistant,
					Content:      query,
					FunctionCall: &FunctionCallObj,
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

// Call the ChatGPT API with passed string and using a prompt and passing a JSON schema
func callChatGPTFunctionCallNoPrompt(query string, prompt string, model string, funcCall []byte) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	var FunctionCallObj openai.FunctionCall
	FunctionCallObj.Name = "functioncall"
	FunctionCallObj.Arguments = string(funcCall)

	resp, err := client.CreateChatCompletion(
		context.Background(),

		// https://platform.openai.com/docs/guides/chat/chat-vs-completions

		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:         openai.ChatMessageRoleAssistant,
					Content:      query,
					FunctionCall: &FunctionCallObj,
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

// Call the GPT Completions API UNUSED CURRENTLY IT SEEMS
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

// Handle a chat interaction with the GPT API
func gptChat(model string, fileChat bool, proglanguage string, file ...string) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	messages := make([]openai.ChatCompletionMessage, 0)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Conversation")
	fmt.Println("---------------------")

	// Create the file to save the chats in
	saveDir := os.Getenv("THYME_QUERY_LOGGING_DIR")
	savefilename, _, _ := makeSaveNameAndStamps(saveDir, "chat")

	// Make the spinner channel so we can tell when its done
	spinningComplete := make(chan bool)

	chatCount := 0

	// If we're reading from a file, read it and send it to the API

	for {

		qs := QuerySave{}

		// If we are filechatting
		if fileChat && chatCount == 0 {

			// Start the spinner
			go spinner(spinningComplete)

			prompt := "Hello! We would like to ask some questions about this file, please:"
			text := readFileToString(file[0])
			text = strings.Replace(text, "\n", "", -1)
			sendtext := fmt.Sprintf("%s %s", prompt, text)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: sendtext,
			})

			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model:    model,
					Messages: messages,
				},
			)

			if err != nil {
				spinningComplete <- true
				fmt.Printf("ChatCompletion error: %v\n", err)
				continue
			}

			// We just save the filename so we dont just create a copy of a
			// Giant file
			qs.Query = fmt.Sprintf("%s %s", prompt, file[0])
			qs.Answer = resp.Choices[0].Message.Content

			saveChat(qs, savefilename)

			chatCount++
			spinningComplete <- true
			continue

		}

		prettyPrintChatArrow("-> ")
		text, _ := reader.ReadString('\n')
		fmt.Println("---")
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})

		// Start the spinner
		go spinner(spinningComplete)

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    model,
				Messages: messages,
			},
		)

		if err != nil {
			spinningComplete <- true
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}

		content := resp.Choices[0].Message.Content
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})

		spinningComplete <- true

		// Save the chat file
		qs.Query = text
		qs.Answer = content
		saveChat(qs, savefilename)

		content = formatCodeBlocksInMarkdown(content, proglanguage)

		typeWriterPrint(content+"\n", false)

		chatCount++
	}
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

	filename, _, formattingTimeStamp := makeSaveNameAndStamps(saveDir, "query")

	var t []byte
	var p []byte
	var q []byte
	var a []byte

	t, err := json.Marshal(formattingTimeStamp)
	if err != nil {
		t = []byte(formattingTimeStamp)
	}

	p, err = json.Marshal(qs.Prompt)
	if err != nil {
		p = []byte(qs.Prompt)
	}

	q, err = json.Marshal(qs.Query)
	if err != nil {
		q = []byte(qs.Query)
	}

	a, err = json.Marshal(qs.Answer)
	if err != nil {
		a = []byte(qs.Answer)
	}

	// Write the file as a json sting string
	// {'timestampe': '...', 'prompt': '...', 'query': '...', 'answer': ...}
	// Filename is YYYY-MM-DD-HH-mm-SS-query.json

	fileData := fmt.Sprintf(`{"timestamp": %s, "prompt": %s, "query": %s, "answer": %s}`,
		string(t),
		string(p),
		string(q),
		string(a))

	// Write the file
	err = ioutil.WriteFile(filename, []byte(fileData), 0644)

	if err != nil {
		fmt.Println(err)
	}
}

/////////////////

// Save the chat information, appended to the file
func saveChat(qs QuerySave, savefile string) {
	file, err := os.OpenFile(savefile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	var q []byte
	var a []byte

	q, err = json.Marshal(qs.Query)
	if err != nil {
		q = []byte(qs.Query)
	}

	a, err = json.Marshal(qs.Answer)
	if err != nil {
		a = []byte(qs.Answer)
	}

	savestr := fmt.Sprintf(`{"query": %s, "answer": %s}`, q, a) + "\n"

	if _, err := file.WriteString(savestr); err != nil { // Append text to file
		fmt.Println(err.Error())
		return
	}

}

// Returns the file save name and timestamps used by saving processes
func makeSaveNameAndStamps(saveDir string, savetype string) (string, string, string) {
	currentTime := time.Now()
	year, month, day := currentTime.Date()
	hour, min, sec := currentTime.Clock()

	formattedTime := fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d",
		year,
		month,
		day,
		hour,
		min,
		sec)

	formattingTimeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		year,
		month,
		day,
		hour,
		min,
		sec)

	filename := fmt.Sprintf("%s/%s-%s.json", saveDir, formattedTime, savetype)

	// If we are saving a chat, make a JSONL file instead
	if savetype == "chat" {
		filename = fmt.Sprintf("%s/%s-%s.jsonl", saveDir, formattedTime, savetype)
	}

	return filename, formattedTime, formattingTimeStamp
}
