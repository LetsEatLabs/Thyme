package main

import "os"

type KagiResponse struct {
	Output string
	Tokens int
}

type KagiRequest struct {
	Engine string
	Input  string // Can be text, a URL, or a filename
	Type   string // url, text, or file
}

func makeURLSummaryRequest(url string, kagi KagiRequest) KagiResponse {
	kagiKey := os.Getenv("KAGI_API_KEY")
}

func makeTextSummaryRequest(url string, kagi KagiRequest) KagiResponse {
	//kagiKey := os.Getenv("KAGI_API_KEY")

	response := KagiResponse{Output: "Not Implemented", Tokens: 0}
	return response
}

func makeFileSummaryRequest(url string, kagi KagiRequest) KagiResponse {
	//kagiKey := os.Getenv("KAGI_API_KEY")

	response := KagiResponse{Output: "Not Implemented", Tokens: 0}
	return response
}
