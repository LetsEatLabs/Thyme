package main

////////////

type Prompt struct {
	Name        string
	Text        string
	Description string
	Examples    map[int]PromptExample // Number them in order to be passed
}

type PromptExample struct {
	Name string //example_assistant, example_user
	Text string
}

////////////

// All prompts are defined in this function initPrompts()
// We may end up moving prompts to JSON files in an embed.FS object if they get too large
// But for now we do not see that being required.

////////////

// Initialize the map of prompt names to prompt structs
// Defined prompts and examples in here
func initPrompts() map[string]Prompt {
	prompts := make(map[string]Prompt)

	// Define prompts here
	prompts["listify"] = Prompt{
		Name:        "listify",
		Text:        "Return a numbered list of actions items from the following text",
		Description: "This prompt takes a block of text and returns a numbered list of action items.",
	}

	prompts["active_voice"] = Prompt{
		Name:        "active_voice",
		Text:        "Write the following text in a more active voice",
		Description: "This prompt takes a block of text and returns a version of the text in a more active voice.",
	}

	prompts["gender-neutral"] = Prompt{
		Name:        "gender-neutral",
		Text:        "Please examine the following text for gendered language and repeat the text back to be but with a non-gendered alternative language, without any explanations:",
		Description: "This prompt takes a block of text and returns a version of the text with gendered language replaced with non-gendered alternatives.",
	}

	return prompts
}

////////////
