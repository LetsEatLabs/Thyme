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
    -p <prompt>         The prompt to use for the GPT request
    -a <ask-question>   Ask a question and get a response (cannot be used with any other flags)
    -h/--help           Display this help message
    -quiet              Will omit the spinner and typewriter. 
    -l                  List all available prompts (-p) and their descriptions. Will exit.
    -text               Pass text to the prompt instead of a file. Used after -p.
                        Anything after is passed. Example: thyme -p active_voice --text "blah"
    -model              
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

    animationFlagVal := flag.Bool("quiet", false, "Will omit the spinner and typewriter.")
    listFlag := flag.Bool("l", false, "List all available prompts (-p) and their descriptions. Will exit.")
    questionFlag := flag.String("a", "", "Ask a question and get a response")
    promptFlag := flag.String("p", "", "The prompt to use for the GPT request: thyme -p active_voice my_blog_post.txt")
    customPromptFlag := flag.String("c", "", "Pass a custom prompt to the GPT request. Cannot be used with -p.")
    textFlag := flag.String("text", "", "Pass text to the prompt instead of a file. Used after -p. Anything after is passed. Example: thyme -p active_voice --text \"blah\"")
    modelFlag := flag.String("model", "chatgpt", "The model to use for the GPT request [chatgpt, gpt4]. Default is chatgpt")
    flag.Parse()

    // A map of string names to our models
    models := map[string]string{
        "chatgpt": openai.GPT3Dot5Turbo,
        "gpt4":    openai.GPT4,
    }

    // If the user passed -l, list the available prompts and exit
    if *listFlag == true {
        listAvailablePrompts()
        os.Exit(0)
    }

    // If the user passed _both_ -c and -p we need to tell them this is not supported
    if *customPromptFlag != "" && *promptFlag != "" {
        fmt.Println("You cannot use both -c and -p. Please use one or the other.")
        os.Exit(1)
    }

    // Make the spinner channel so we can tell when its done
    spinningComplete := make(chan bool)

    if *animationFlagVal == false {
        go spinner(spinningComplete)
    }

    // Get the value after the -p flag
    prompt := promptFlag
    var request string
    var chosenPrompt string

    // If the user passed -text then we want to use the text after the flag
    if *textFlag == "" && *questionFlag == "" {
        request = readFileToString(flag.Args()[0])
    } else if *textFlag == "" {
        request = *questionFlag
    } else {
        request = *textFlag
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
            response = callChatGPTNoPrompt(request, models[*modelFlag])
        } else {
            response = callChatGPT(request, chosenPrompt, models[*modelFlag])
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
            typeWriterPrint(cleanResponse)
        } else {
            fmt.Println(cleanResponse)
        }

        os.Exit(0)
    }

    response := callChatGPT(request, chosenPrompt, models[*modelFlag])

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
        typeWriterPrint(cleanResponse)
    } else {
        fmt.Println(cleanResponse)
    }

}
