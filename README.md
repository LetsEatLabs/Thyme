# Thyme

A CLI helper for interacting with the OpenAI API. Capable of sending files using built-in prompts, examples (coming soon), and prompt-chains (coming soon). Also capable of simple direct queries.

Useful as an inter-process tool, to use as a "function call" from your other applications. Or just use it as a CLI tool, you do you.


## Installation

```bash
~ $: export OPENAI_API_KEY="<your_key_here>"
~ $: go build -o thyme *.go
~ $: ./thyme --help                                                                                                                                        
Usage of thyme:
  -a string
        Ask a question and get a response
  -c string
        Pass a custom prompt to the GPT request. Cannot be used with -p.
  -l    List all available prompts (-p) and their descriptions. Will exit.
  -model string
        The model to use for the GPT request [chatgpt, gpt4]. Default is chatgpt (default "chatgpt")
  -p string
        The prompt to use for the GPT request: thyme -p active_voice my_blog_post.txt
  -quiet
        Will omit the spinner and typewriter.
  -text string
        Pass text to the prompt instead of a file. Used after -p. Anything after is passed. Example: thyme -p active_voice --text "blah"

```

You can enable saving the queries and results as JSON with the following environment variables:

`THYME_QUERY_LOGGING_DIR='<full_path_to_dir>'`
`THYME_QUERY_LOGGING='true'`

If anything but 'true' is set for `THYME_QUERY_LOGGING` then it will not be logged.



## Downloads
MacOS build: [thyme](https://static.letseatlabs.com/bin/thyme/macos/thyme)
Linux build: [thyme](https://static.letseatlabs.com/bin/thyme/linux/thyme)


## Examples

```bash
~ $: thyme -c "Please take the following text and return a JSON object of the different word types such as verb, nouns, etc. Please do not explain anything." -text "Today I went to the park and tomorrow I need to go to the zoo. After the store today I will eat a hamburger" 
{                   
  "noun": [
    "park",
    "tomorrow",
    "zoo",
    "store",
    "hamburger"
  ],
  "verb": [
    "went",
    "need",
    "will",
    "eat"
  ],
  "adjective": [
    "Today"
  ]
}
```

```bash
~ $: thyme -p listify test.txt
1. Go to the store and buy eggs and milk.
2. Stop by the gasoline store and buy tires.
3. Purchase a whole salmon and eat it.
```

```bash
~ $: thyme -model gpt4 -p summarize-text README.md
Thyme is a CLI helper for interacting with the OpenAI API, designed for sending files and simple direct queries. It can be used as an inter-process tool or a standalone CLI tool. The installation process requires the users to export their OpenAI API key and build the application. Examples provided demonstrate how to use Thyme for various purposes, such as asking a question and processing a text file to list the content.

Notable points:
1. CLI helper for OpenAI API interaction.
2. Supports built-in prompts, examples, and prompt-chains.
3. Can be used as an inter-process tool or standalone CLI tool.
4. Requires OpenAI API key for installation.
5. Provides usage examples and instructions in the documentation.
```

## ToDo
- [x] Custom prompts used with files
- [x] Add prompt to summarize body of text
