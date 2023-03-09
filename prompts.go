package main

////////////

type Prompt struct {
	Name     string
	Text     string
	Examples map[string]string
}

////////////

// Initialize the map of prompt names to prompt structs
func initPrompts() map[string]Prompt {
	prompts := make(map[string]Prompt)

	prompts["listify"] = Prompt{
		Name: "listify",
		Text: "Return a numbered list of actions items from the following text",
	}

	return prompts
}

////////////
