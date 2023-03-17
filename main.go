package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
)

// Help message to display when the user asks for help or
// fails to pass any arguments
func helpMessage() {
    helpStr := `
Usage: thyme <flags> <input file>

Flags:
    -p <prompt>         The prompt to use for the GPT request
    -a <ask-question>   Ask a question and get a response (cannot be used with any other flags)
    -h                  Display this help message
    -quiet              Will omit the spinner and typewriter. 
    -l                  List all available prompts (-p) and their descriptions. Will exit.
    -text               Pass text to the prompt instead of a file. Used after -p.
                        Anything after is passed. Example: thyme -p active_voice --text "blah"
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

    // Parse arguements and load the prompts struct
    //arguments, _ := parseArgs()
    prompts := initPrompts()

    animationFlagVal := flag.Bool("quiet", false, "Will omit the spinner and typewriter.")
    helpFlag := flag.Bool("h", false, "Display this help message")
    listFlag := flag.Bool("l", false, "List all available prompts (-p) and their descriptions. Will exit.")
    questionFlag := flag.String("a", "", "Ask a question and get a response")
    promptFlag := flag.String("p", "", "The prompt to use for the GPT request")
    textFlag := flag.String("text", "", "Pass text to the prompt instead of a file. Used after -p. Anything after is passed. Example: thyme -p active_voice --text \"blah\"")
    flag.Parse()

    // If the user passed -h or --help, display the help message and exit
    if *helpFlag == true {
        helpMessage()
        os.Exit(0)
    }

    // If the user passed -l, list the available prompts and exit
    if *listFlag == true {
        listAvailablePrompts()
        os.Exit(0)
    }

    // Make the spinner channel so we can tell when its done
    spinningComplete := make(chan bool)

    if *animationFlagVal == false {
        go spinner(spinningComplete)
    }

    // -a flag will allow us to just ask a question and get a response
    if *questionFlag != "" {
        var request *string
        var response string

        request = questionFlag
        response = callChatGPTNoPrompt(*request)

        // Tell the spinner we are done
        if *animationFlagVal == false {
            spinningComplete <- true
        }

        cleanResponse := removeLeadingNewLines(response)

        if *animationFlagVal == false {
            typeWriterPrint(cleanResponse)
        } else {
            fmt.Println(cleanResponse)
        }

        os.Exit(0)
    }

    // Get the value after the -p flag if we are not doing a -q request
    prompt := promptFlag
    var request string

    // If the user passed --text then we want to use the text after the flag
    if *textFlag == "" {
        request = readFileToString(flag.Args()[0])
    } else {
        request = *textFlag
    }

    response := callChatGPT(request, prompts[*prompt].Text)

    // Tell the spinner we are done

    if *animationFlagVal == false {
        spinningComplete <- true
    }

    cleanResponse := removeLeadingNewLines(response)

    if *animationFlagVal == false {
        typeWriterPrint(cleanResponse)
    } else {
        fmt.Println(cleanResponse)
    }

}
