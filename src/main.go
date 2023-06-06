package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// Help message to display when the user asks for help or
// fails to pass any arguments
func helpMessage() {
	helpStr := `
Usage: thyme <flags> <input file>

Flags:
  -a string
        Ask a question and get a response
  -c string
        Pass a custom prompt to the GPT request. Cannot be used with -p.
  -chat
        Start a chat session with the GPT model. Must be used with -oa. Can be used with -file to chat about a file.
  -file string
        Pass file to the prompt. Cannot be used with -a.
  -history string
        Review the history of your queries, or a specific one. -history [chat, summary, query, all, <full-path-to-history-file>]
  -ksum string
        Use the Kagi Universal Summarizer API. -ksum [text | url]. Also works with -model
  -ktype string
        Type of summary from the Kagi Universal Summarizer API. -ktype [summary,notes]. 'summary' gives a paragraph, 'notes' gives points.
  -l    List all available prompts (-p) and their descriptions. Will exit.
  -lang string
        The language to format the response syntax for. Omit to 'guess'.
  -model string
        The model to use for the request. OpenAI: [chatgpt, gpt4] Kagi: [agnes, daphne, muriel($$)]. Defaults are chatgpt and agnes.
  -oa
        Use the OpenAI API.
  -p string
        The prompt to use for the GPT request: thyme -p active_voice my_blog_post.txt
  -quiet
        Will omit the spinner, typewriter, and color effects.
          
`

	fmt.Println(helpStr)
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

	// If not arguments passed, display help message and exit
	if len(os.Args) == 1 {
		helpMessage()
		os.Exit(0)
	}

	// If the env argument OPEN_AI_API key does not exist, exit
	// with an error message
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable")
		os.Exit(1)
	}

	// Are we saving queries today?
	sq := os.Getenv("THYME_QUERY_LOGGING")
	saveQueries := false

	if sq == "true" {
		saveQueries = true
	} else {
		saveQueries = false
	}

	// Parse arguements and load the prompts struct
	//arguments, _ := parseArgs()
	prompts := initPrompts()

	animationFlagVal := flag.Bool("quiet", false, "Will omit the spinner, typewriter, and color effects.")
	listFlag := flag.Bool("l", false, "List all available prompts (-p) and their descriptions. Will exit.")
	questionFlag := flag.String("a", "", "Ask a question and get a response")
	promptFlag := flag.String("p", "", "The prompt to use for the GPT request: thyme -p active_voice my_blog_post.txt")
	customPromptFlag := flag.String("c", "", "Pass a custom prompt to the GPT request. Cannot be used with -p.")
	modelFlag := flag.String("model", "", "The model to use for the request. OpenAI: [chatgpt, gpt4] Kagi: [agnes, daphne, muriel($$)]. Defaults are chatgpt and agnes.")
	chatFlag := flag.Bool("chat", false, "Start a chat session with the GPT model. Must be used with -oa. Can be used with -file to chat about a file.")
	kagiFlag := flag.String("ksum", "", "Use the Kagi Universal Summarizer API. -ksum [text | url]. Also works with -model")
	kagiGPTFlag := flag.Bool("kgpt", false, "Use the Kagi FastGPT API. -ksum [query text]. Always defaults to web_search=true")
	kagiTypeFlag := flag.String("ktype", "", "Type of summary from the Kagi Universal Summarizer API. -ktype [summary,notes]. 'summary' gives a paragraph, 'notes' gives points.")
	openAIFlag := flag.Bool("oa", false, "Use the OpenAI API.")
	fileFlag := flag.String("file", "", "Pass file to the prompt. Cannot be used with -a.")
	langFlag := flag.String("lang", "", "The language to format the response syntax for. Omit to 'guess'.")
	historyFlag := flag.String("history", "", "Review the history of your queries, or a specific one. -history [chat, summary, query, all, <full-path-to-history-file>]")
	flag.Parse()

	// A map of string names to our models
	openAIModels := map[string]string{
		"chatgpt": openai.GPT3Dot5Turbo,
		"gpt4":    openai.GPT4,
	}

	kagiModels := map[string]string{
		"agnes":  "agnes",
		"daphne": "daphne",
		"muriel": "muriel",
	}

	// If the user passed -l, list the available prompts and exit
	if *listFlag == true {
		listAvailablePrompts()
		os.Exit(0)
	}

	// If we want to view our history, do that
	if *historyFlag != "" {
		viewHistoryQueries(*historyFlag)
		os.Exit(0)
	}

	// If the user passed _both_ -c and -p we need to tell them this is not supported
	if *customPromptFlag != "" && *promptFlag != "" {
		fmt.Println("You cannot use both -c and -p. Please use one or the other.")
		os.Exit(1)
	}

	// If they passed both a file an a question tell them no
	if *fileFlag != "" && *questionFlag != "" {
		fmt.Println("You cannot use both -file and -a. Please use one or the other.")
		os.Exit(1)
	}

	// Make the spinner channel so we can tell when its done
	spinningComplete := make(chan bool)

	// Handle requests. This is the order flags are checked in. Each one will exit
	// So we only ever use one call per execution.

	// Handle a Kagi API
	if *kagiFlag != "" || *kagiGPTFlag == true {

		// Start the spinner
		if *animationFlagVal == false {
			go spinner(spinningComplete)
		}

		if *questionFlag == "" {
			fmt.Println("Please pass either a URL or text after -a: thyme -ksum -a https://a.com")
			os.Exit(1)
		}

		// We default to agnes, but if the user passes a different one we use that
		var engineChoice string

		if *modelFlag != "" {
			engineChoice = kagiModels[*modelFlag]
		} else {
			engineChoice = "agnes"
		}

		kagi := KagiRequest{
			Engine:      engineChoice,
			Input:       *questionFlag,
			Type:        *kagiFlag,
			SummaryType: *kagiTypeFlag,
		}

		if *kagiGPTFlag != false {
			kagi.Type = "fastgpt"
			kagi.SummaryType = "fastgpt"
			kagi.Engine = "fastgpt"
		}

		response := makeKagiRequest(kagi)

		// Tell the spinner we are done and print the response
		if *animationFlagVal == false {
			spinningComplete <- true
			typeWriterPrint(response.Data.Output, false)
			os.Exit(0)
		}

		fmt.Println(response.Data.Output)
		os.Exit(0)

	}

	// Handle an OpenAI Request
	if *openAIFlag == true {

		// We default to chatgpt, but if the user passes a different one we use that
		var engineChoice string

		if *modelFlag != "" {
			engineChoice = *modelFlag
		} else {
			engineChoice = "chatgpt"
		}

		// If the user wishes to chat, lets do that
		if *chatFlag == true {

			// If the user wants to chat about a file
			if *fileFlag != "" {
				gptChat(openAIModels[engineChoice], true, *langFlag, *fileFlag)
				os.Exit(0)
			}

			gptChat(openAIModels[engineChoice], false, *langFlag)
			os.Exit(0)
		}

		// Enable the spinner if it is not disabled
		if *animationFlagVal == false {
			go spinner(spinningComplete)
		}

		// Get the value after the -p flag
		prompt := promptFlag
		var request string
		var chosenPrompt string

		// If the user passed -a then we want to use the text after the flag
		if *fileFlag != "" {
			request = readFileToString(*fileFlag)
		} else if *questionFlag != "" {
			request = *questionFlag
		}

		// Load the prompt from the list
		// If we passed a -c flag, replace the prompt with the custom text

		if *customPromptFlag != "" {
			chosenPrompt = *customPromptFlag
		} else {
			chosenPrompt = prompts[*prompt].Text
		}

		// -a flag will allow us to just ask a question and get a response
		// Exit after we display the response
		if *questionFlag != "" {

			var response string

			if *promptFlag != "" && *customPromptFlag != "" {
				response = callChatGPTNoPrompt(request, openAIModels[engineChoice])
			} else {
				response = callChatGPT(request, chosenPrompt, openAIModels[engineChoice])
			}

			// Tell the spinner we are done
			if *animationFlagVal == false {
				spinningComplete <- true
			}

			cleanResponse := removeLeadingNewLines(response)

			// Save query before we display it incase user ctrl-c's and its still logged
			qs := QuerySave{
				Query:  request,
				Prompt: chosenPrompt,
				Answer: cleanResponse,
			}

			if saveQueries {
				saveGPT(qs)
			}

			if *animationFlagVal == false {
				cleanResponse = formatCodeBlocksInMarkdown(cleanResponse, *langFlag)
				typeWriterPrint(cleanResponse, true)
			} else {
				fmt.Println(cleanResponse)
			}

			os.Exit(0)
		}

		response := callChatGPT(request, chosenPrompt, openAIModels[engineChoice])

		// Tell the spinner we are done

		if *animationFlagVal == false {
			spinningComplete <- true
		}

		cleanResponse := removeLeadingNewLines(response)

		// Save query before we display it incase user ctrl-c's and its still logged
		qs := QuerySave{
			Query:  request,
			Prompt: chosenPrompt,
			Answer: cleanResponse,
		}

		if saveQueries {
			saveGPT(qs)
		}

		if *animationFlagVal == false {
			cleanResponse = formatCodeBlocksInMarkdown(cleanResponse, *langFlag)
			typeWriterPrint(cleanResponse, true)
		} else {
			fmt.Println(cleanResponse)
		}

	}
}
