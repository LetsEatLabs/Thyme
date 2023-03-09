package main

////////////

type Prompt struct {
	Name     string
	Text     string
	Examples map[string]PromptExample
}

type PromptExample struct {
	Name string //example_assistant, example_user
	Text string
}

////////////

// Initialize the map of prompt names to prompt structs
// Defined prompts and examples in here
func initPrompts() map[string]Prompt {
	prompts := make(map[string]Prompt)

	// Define prompts here
	prompts["listify"] = Prompt{
		Name: "listify",
		Text: "Return a numbered list of actions items from the following text",
	}

	prompts["active_voice"] = Prompt{
		Name: "active_voice",
		Text: "Write the following text in a more active voice",
	}

	return prompts
}

////////////
