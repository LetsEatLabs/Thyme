package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

// Read the prompts.json file and return it as a map
func readPromptsConfig() map[string]interface{} {
	var result map[string]interface{}
	promptsConfigString := readFileToString("prompts.json")
	err := json.Unmarshal([]byte(promptsConfigString), &result)

	if err != nil {
		fmt.Println(err)
	}

	return result
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
	promptsConfig := readPromptsConfig()

	prompt := arguments["-p"]

	request := readFileToString(arguments["input"])
	response := callChatGPT(request, promptsConfig[prompt].(string))
	fmt.Println(response)
}
