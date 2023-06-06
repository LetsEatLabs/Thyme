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
	kagiFastGPTEndpoint    = "https://kagi.com/api/v0/fastgpt"
	kagiKey                = os.Getenv("KAGI_API_KEY")
	kagiAuthHeader         = fmt.Sprintf("Bot %s", kagiKey)
)

type KagiResponse struct {
	Data struct {
		Output     string       `json:"output"`
		Tokens     int          `json:"tokens"`
		References []KagiSource `json:"references"`
	} `json:"data"`

	Meta struct {
		Id   string `json:"id"`
		Node string `json:"node"`
		Ms   int    `json:"ms"`
	} `json:"meta"`
}

type KagiRequest struct {
	Engine      string
	Input       string // Can be text, a URL, or a filename
	Type        string // url, text, or file
	SummaryType string // summary or notes (points)
}

type KagiSource struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	URL     string `json:"url"`
}

func makeKagiRequest(kagi KagiRequest) KagiResponse {

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

	// Get the summary type. By default anything but "notes" or "fastgpt" is a "summary"
	if kagi.SummaryType != "notes" && kagi.SummaryType != "fastgpt" {
		kagi.SummaryType = "summary"
	}

	// It is actually called takeaway but notes is shorter to type
	if kagi.SummaryType == "notes" {
		kagi.SummaryType = "takeaway"
	}

	cleanSumType, err := json.Marshal(kagi.SummaryType)
	if err != nil {
		fmt.Println("Error marshalling sumType: ", err)
		os.Exit(1)
	}

	// Initialize the request string
	var request string

	if kagi.Type == "fastgpt" {
		request = fmt.Sprintf(`{"query": %s, "web_search": true, "cache": true}`, cleanInput)
	} else {
		request = fmt.Sprintf(`{"%s": %s, "engine": %s, "summary_type": %s}`,
			kagi.Type, cleanInput, cleanEngine, cleanSumType)
	}

	brequest := []byte(request) // Bytes so we can send it over the wire

	// Create a new request with custom headers and JSON payload
	var usingEndpoint string

	if kagi.Type == "fastgpt" {
		usingEndpoint = kagiFastGPTEndpoint
	} else {
		usingEndpoint = kagiURLSummaryEndpoint
	}

	req, err := http.NewRequest("POST", usingEndpoint, bytes.NewBuffer(brequest))
	if err != nil {
		fmt.Println(err)
	}

	// Apply the headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Convert response body to JSON
	var response KagiResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
	}

	if kagi.Type == "fastgpt" {
		response.Data.Output += "\n\n" + kagiSourcesToString(response.Data.References)
	}

	if os.Getenv("THYME_QUERY_LOGGING") == "true" {
		saveKagiSummary(response, kagi)
	}

	return response
}

func kagiSourcesToString(sources []KagiSource) string {
	var output string

	output += "Sources:\n----------\n"

	for a, source := range sources {
		output += fmt.Sprintf("[%d] %s\n%s\n%s\n\n", a+1, source.Title, source.URL, source.Snippet)
	}

	return output
}

func saveKagiSummary(response KagiResponse, request KagiRequest) {
	directory := os.Getenv("THYME_QUERY_KAGI_LOGGING_DIR")
	fileloc, _, _ := makeSaveNameAndStamps(directory, "summary")

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
