package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ///
func readPromptsConfig() map[string]interface{} {
	var result map[string]interface{}
	promptsConfigString := readFileToString("prompts.json")
	err := json.Unmarshal([]byte(promptsConfigString), &result)

	if err != nil {
		fmt.Println(err)
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
	promptsConfig := readPromptsConfig()
	request := readFileToString("test.txt")
	response := callChatGPT(request, promptsConfig["listify"].(string))
	fmt.Println(response)
}
