package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Help message to display when the user asks for help or
// fails to pass any arguments
func helpMessage() {
	helpStr := `
Usage: chatgpt <flags> <input file>

Flags:
    -p <prompt>     The prompt to use for the chatbot
    -h (--help)     Display this help message
`

	fmt.Println(helpStr)
}

// Parses the command line arguments and returns them in a map
// of flag:argument, if there are an uneven number of arguments then we
// name the final argument "input" and add it to the map.
// If the user asks for help or fails to pass any
// Display the help message and exit.
func parseArgs() map[string]string {

	if len(os.Args) == 1 {
		helpMessage()
		os.Exit(0)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		helpMessage()
		os.Exit(0)
	}

	args := os.Args[1:]
	result := make(map[string]string)
	for i := 0; i < len(args); i += 2 {

		if i+2 > len(args) {
			result["input"] = args[i]
		} else {
			result[args[i]] = args[i+1]
		}

	}
	return result
}

/////

// Read a file and return its contents as a string
func readFileToString(filename string) string {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
	}

	return string(file)
}

func main() {
	arguments := parseArgs()
	prompts := initPrompts()

	// -q flag will allow us to just ask a question and get a response
	_, ok := arguments["-q"]
	if ok {
		request := strings.Join(os.Args[2:], " ")
		response := callChatGPTNoPrompt(request)
		cleanResponse := removeLeadingNewLines(response)
		typeWriterPrint(cleanResponse)
		os.Exit(0)
	}

	prompt := arguments["-p"]

	request := readFileToString(arguments["input"])
	response := callChatGPT(request, prompts[prompt].Text)
	cleanResponse := removeLeadingNewLines(response)
	typeWriterPrint(cleanResponse)

}
