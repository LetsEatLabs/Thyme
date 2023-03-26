package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	kagiURLSummaryEndpoint = "https://kagi.com/api/v0/summarize"
	kagiKey                = os.Getenv("KAGI_API_KEY")
	kagiAuthHeader         = fmt.Sprintf("Bot %s", kagiKey)
)

type KagiResponse struct {
	Data struct {
		Output string `json:"output"`
		Tokens int    `json:"tokens"`
	} `json:"data"`

	Meta struct {
		Id   string `json:"id"`
		Node string `json:"node"`
		Ms   int    `json:"ms"`
	} `json:"meta"`
}

type KagiRequest struct {
	Engine string
	Input  string // Can be text, a URL, or a filename
	Type   string // url, text, or file
}

func makeSummaryRequest(kagi KagiRequest) KagiResponse {

	// Sanitize our input to make a JSON string
	cleanInput, err := json.Marshal(kagi.Input)
	if err != nil {
		fmt.Println("Error marshalling input: ", err)
		os.Exit(1)
	}

	cleanEngine, err := json.Marshal(kagi.Engine)
	if err != nil {
		fmt.Println("Error marshalling engine: ", err)
		os.Exit(1)
	}

	// Set custom headers
	headers := map[string]string{
		"Authorization": kagiAuthHeader,
		"Content-Type":  "application/json",
	}

	request := fmt.Sprintf(`{"%s": %s, "engine": %s}`, kagi.Type, cleanInput, cleanEngine)
	brequest := []byte(request) // Bytes so we can send it over the wire

	// Create a new request with custom headers and JSON payload
	req, err := http.NewRequest("POST", kagiURLSummaryEndpoint, bytes.NewBuffer(brequest))
	if err != nil {
		panic(err)
	}

	// Apply the headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Convert response body to JSON
	var response KagiResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	if os.Getenv("THYME_QUERY_LOGGING") == "true" {
		saveKagiSummary(response, kagi)
	}

	return response
}

// func makeFileSummaryRequest(kagi KagiRequest) KagiResponse {
// 	//kagiKey := os.Getenv("KAGI_API_KEY")

// 	response := KagiResponse{Output: "Not Implemented", Tokens: 0}
// 	return response
// }

func saveKagiSummary(response KagiResponse, request KagiRequest) {
	directory := os.Getenv("THYME_QUERY_KAGI_LOGGING_DIR")
	fileloc, _, _ := makeSaveNameAndStamps(directory)

	var q []byte
	var a []byte

	q, err := json.Marshal(request.Input)
	if err != nil {
		q = []byte(request.Input)
	}

	a, err = json.Marshal(response.Data.Output)
	if err != nil {
		a = []byte(response.Data.Output)
	}

	fileData := fmt.Sprintf(`{"query": %s, "answer": %s}`, string(q), string(a))

	// Write the file
	err = ioutil.WriteFile(fileloc, []byte(fileData), 0644)

	if err != nil {
		fmt.Println(err)
	}

}
