# Thyme

A CLI helper for interacting with multiple AI/LLM APIs. Capable of sending files using built-in prompts, examples (coming soon), and prompt-chains (coming soon). Also capable of simple direct queries and sending files.

```bash
Usage of thyme:
  -a string
        Ask a question and get a response
  -c string
        Pass a custom prompt to the GPT request. Cannot be used with -p.
  -chat
        Start a chat session with the GPT model
  -file string
        Pass file to the prompt. Cannot be used with -a.
  -ksum
        Use the Kagi Universal Summarizer API. Pass the URL to summarize after -a.
  -l    List all available prompts (-p) and their descriptions. Will exit.
  -model string
        The model to use for the GPT request if you wish to not use the default
  -p string
        The prompt to use for the GPT request: thyme -p active_voice my_blog_post.txt
  -quiet
        Will omit the spinner and typewriter.

```


## Installation

You can enable saving the queries and results as JSON with the following environment variables:

- `THYME_QUERY_LOGGING_DIR='<full_path_to_dir>'`
- `THYME_QUERY_LOGGING='true'`

The program accesses the following environment variables

| Variable | Use | Example | Required | 
| --- | --- | --- | --- |
| `OPENAI_API_KEY` | The OpenAI API key | `sk-1234567890` | Yes |
| `KAGI_API_KEY` | The Kagi API key | `AAA_Keysomething12389asd` | Yes |
| `THYME_QUERY_LOGGING` | Whether to log the queries and results | `true` | No |
| `THYME_QUERY_LOGGING_DIR` | The directory to save the query logs to. Currenly only supports OpenAI | `/home/user/.thyme/logs` | No |

If anything but 'true' is set for `THYME_QUERY_LOGGING` then it will not be logged.

The models we currently support interacting with are as follows: 

| Model | Provider | Default |
| --- | --- | --- |
| `chatgpt` | OpenAI | Yes |
| `gpt4` | OpenAI | No |
| `agnes` | Kagi | Yes |
| `daphne` | Kagi | No |

## Use

### Interprocess

The most common use cases for this application are likely to be quick answers to questions or the chat interface. But this was also designed to be used between processes, so you can use it as a "function call" from yosur other applications. By default the application has some animations and formatting, but if you pass the `-quiet` flag then this is omitted and you are simply returned the response.

### Built-in Prompts

To view the list of current built in prompts, please use `thyme -l`.

### Chat

To chat with any of the Open AI models, you can use the `-chat` flag.

### Summarize large bodies of text

You can utilize the Kagi Universal Summarizer API to summarize large bodies of text with `-ksum`. Kagi currently only supports URLs and raw text right now, but they plan to support file upload in the future.


## Examples

Prompts and queries

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

Built-in prompts

```bash
~ $: thyme -p listify test.txt
1. Go to the store and buy eggs and milk.
2. Stop by the gasoline store and buy tires.
3. Purchase a whole salmon and eat it.
```

Alternate-model support

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

Interactive Chat Sessions (supports ChatGPT and GPT-4)

```bash
~ $: thyme -chat            
Conversation
---------------------
-> How many moons does saturn have?
---


Saturn has 82 moons, the largest of which is called Titan.

-> Why Titan?
---
Titan is the largest of Saturn's moons, and it was named after the Titans of Greek mythology, which were powerful giants who were the ancestors of the gods. The name Titan is very appropriate for this moon, as it is the only known moon in the solar system to have a thick atmosphere, with clouds, rain, lakes, and rivers. Titan is also the second-largest moon in the solar system, after Jupiter's moon Ganymede.

-> 
```

## ToDo
- [x] Custom prompts used with files
- [x] Add prompt to summarize body of text
