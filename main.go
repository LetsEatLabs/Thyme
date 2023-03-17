package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

// Help message to display when the user asks for help or
// fails to pass any arguments
func helpMessage() {
    helpStr := `
Usage: thyme <flags> <input file>

Flags:
    -p <prompt>         The prompt to use for the GPT request
    -a <ask-question>       Ask a question and get a response (cannot be used with any other flags)
    -h (--help)         Display this help message
    --animation false   Will omit the spinner and typewriter. 
                        Flag _must_ come first when used with the -q flag. 
                        "false" must be passed. 
    -l                  List all available prompts (-p) and their descriptions. Will exit.
    --text              Pass text to the prompt instead of a file. Used after -p.
                        Anything after is passed. Example: thyme -p active_voice --text "blah"
`

    fmt.Println(helpStr)
}

// Parses the command line arguments and returns them in a map
// of flag:argument, if there are an uneven number of arguments then we
// name the final argument "input" and add it to the map.
// If the user asks for help or fails to pass any
// Display the help message and exit.
// Will return true if the user passed the --text flag
func parseArgs() (map[string]string, bool) {

    if len(os.Args) == 1 {
        helpMessage()
        os.Exit(0)
    }

    if os.Args[1] == "-h" || os.Args[1] == "--help" {
        helpMessage()
        os.Exit(0)
    }

    // If the user wants to see the prompts, print them and exit
    if os.Args[1] == "-l" {
        listAvailablePrompts()
        os.Exit(0)
    }

    args := os.Args[1:]

    result := make(map[string]string)
    for i := 0; i < len(args); i += 2 {

        if i+2 > len(args) {
            result["input"] = args[i]
        } else {

            // If we pass --text, then we want everything after this to be the input
            if args[i] != "--text" {
                result[args[i]] = args[i+1]
            } else {
                result["input"] = strings.Join(args[i+1:], " ")
                //fmt.Println(result)
                return result, true
            }

        }

    }
    //fmt.Println(result)
    return result, false
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

    // If the env argument OPEN_AI_API key does not exist, exit
    // with an error message
    if os.Getenv("OPENAI_API_KEY") == "" {
        fmt.Println("Please set the OPENAI_API_KEY environment variable")
        os.Exit(1)
    }

    // Parse arguements and load the prompts struct
    arguments, cli := parseArgs()
    prompts := initPrompts()

    animationFlagVal := flag.Bool("quiet", false, "Will omit the spinner and typewriter.")
    questionFlag := flag.String("a", "", "Ask a question and get a response")
    promptFlag := flag.String("p", "", "The prompt to use for the GPT request")
    flag.Parse()

    // Make the spinner channel so we can tell when its done
    spinningComplete := make(chan bool)

    if *animationFlagVal == false {
        go spinner(spinningComplete)
    }

    // -a flag will allow us to just ask a question and get a response
    if questionFlag != nil {
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
    if cli != true {
        request = readFileToString(arguments["input"])
    } else {
        request = arguments["input"]
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
