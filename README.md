# Thyme

![Thyme.png](Thyme.png)

A CLI helper for interacting with multiple AI/LLM APIs. Capable of sending files using built-in prompts, examples (coming soon), and prompt-chains (coming soon). Also capable of simple direct queries and sending files.

```bash
Usage of thyme:
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
  -json string
      Give a json schema file to send as a FunctionCall to get a structured response. Currently only works with a prompt (-p or -c). -json <full-path-to-json-schema-file>
  -kgpt
      Use the Kagi FastGPT API. -ksum [query text]. Always defaults to web_search=true
  -ksum string
      Use the Kagi Universal Summarizer API. -ksum [text | url]. Also works with -model
  -ktype string
      Type of summary from the Kagi Universal Summarizer API. -ktype [summary,notes]. 'summary' gives a paragraph, 'notes' gives points.
  -l  List all available prompts (-p) and their descriptions. Will exit.
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

```


## Installation

You can enable saving the queries and results as JSON with the following environment variables:

- `THYME_QUERY_LOGGING_DIR='<full_path_to_dir>'`
- `THYME_QUERY_LOGGING='true'`
- `THYME_QUERY_KAGI_LOGGING_DIR='<full_path_to_dir>'`

The program accesses the following environment variables

| Variable | Use | Example | Required | 
| --- | --- | --- | --- |
| `OPENAI_API_KEY` | The OpenAI API key | `sk-1234567890` | Yes |
| `KAGI_API_KEY` | The Kagi API key | `AAA_Keysomething12389asd` | Yes |
| `THYME_QUERY_LOGGING` | Whether to log the queries and results | `true` | No |
| `THYME_QUERY_LOGGING_DIR` | The directory to save the query logs to. Currenly only supports OpenAI | `/home/user/.thyme/logs` | No |

If anything but 'true' is set for `THYME_QUERY_LOGGING` then it will not be logged.

Now build the program:

```bash
~ $: cd src
~ $: go build -o thyme .
~ $: ./thyme -oa -chat
```

The models we currently support interacting with are as follows: 

| Model | Provider | Default |
| --- | --- | --- |
| `chatgpt` | OpenAI | Yes |
| `gpt4` | OpenAI | No |
| `agnes` | Kagi | Yes |
| `daphne` | Kagi | No |

We also support other features of the Kagi API:

| Feature | Use |
| --- | --- |
| `ksum` | Summarize large bodies of text or a URL |
| `kgpt` | Kagi's FastGPT with web search capabilities. |



## Use

### Interprocess

The most common use cases for this application are likely to be quick answers to questions or the chat interface. But this was also designed to be used between processes, so you can use it as a "function call" from your other applications. By default the application has some animations and formatting, but if you pass the `-quiet` flag then this is omitted and you are simply returned the response.

### Built-in Prompts

To view the list of current built in prompts, please use `thyme -l`.

### Chat

To chat with any of the Open AI models, you can use the `-chat` flag.

### Summarize large bodies of text

You can utilize the Kagi Universal Summarizer API to summarize large bodies of text with `-ksum`. Kagi currently only supports URLs and raw text right now, but they plan to support file upload in the future.


## Examples

Prompts and queries

```bash
~ $: thyme -oa chatgpt -c "Please take the following text and return a JSON object of the different word types such as verb, nouns, etc. Please do not explain anything." -text "Today I went to the park and tomorrow I need to go to the zoo. After the store today I will eat a hamburger" 
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
~ $: thyme -oa chatgpt -p listify test.txt
1. Go to the store and buy eggs and milk.
2. Stop by the gasoline store and buy tires.
3. Purchase a whole salmon and eat it.
```

Alternate-model support

```bash
~ $: thyme -oa -model gpt4 -a "Vic20 BASIC program that displays a blue square and rotates it, no explanation"
10 PRINT CHR$(147) : REM CLEAR SCREEN
20 REM SET UP COLOR MEMORY
30 POKE 36879, 6 : REM SET BORDER COLOR TO BLUE
40 POKE 36878, 14 : REM SET BACKGROUND COLOR TO BLUE
50 POKE 36875, 0 : REM SET SCREEN MEMORY TO 0
60 POKE 648, 0 : REM SET CHARACTER COLOR TO BLACK
70 REM DRAW SQUARE
80 FOR I=10 TO 13
90 Y=I*22
100 FOR X=Y TO Y+3
110 POKE 1024+X, 160
120 POKE 55296+X, 1
130 NEXT X
140 NEXT I
150 REM ROTATE SQUARE
160 FOR ITER=1 TO 100
170 POKE 36875, 0+16*(ITER MOD 2) : REM TOGGLE SCREEN MEMORY BETWEEN 0 AND 16
180 POKE 36877, 204+4*(ITER MOD 4) : REM SWITCH TO BUFFERED CHARACTER SETS
190 FOR DELAY=1 TO 200 : NEXT DELAY : REM PAUSE FOR A WHILE
200 NEXT ITER
210 GOTO 150
To run this program on a VIC-20, you'll need to replace the POKE and locations using the VIC-20 equivalents. In any case, the program will display a rotating blue square without giving any explanation.
```

Interactive Chat Sessions (supports ChatGPT and GPT-4)

```bash
~ $: thyme -oa -chat            
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

Using the Kagi text summary API to summarize the [599 word sentence](https://nathanbrixius.wordpress.com/2013/10/30/the-five-longest-proust-sentences/) from Swann's Way by Proust

```bash
~ $: thyme -ksum text -a "[OMITTED FOR BREVITY]"
The document describes the author's experience of revisiting the various rooms he had slept in throughout his life in a long dream. He describes the different types of rooms he had slept in, such as rooms in winter where he would feel warm and cozy, and rooms in summer where he would feel a part of the warm evening. He also describes specific rooms, such as the Louis XVI room, which was so cheerful that he could never feel unhappy in it, and a little room with a high ceiling and mahogany walls, where he felt anxious and uncomfortable. The author describes how his mind would elongate itself upwards to take on the exact shape of the room, and how he would spend anxious nights until he became accustomed to the room and its surroundings. The author also describes how his perception of the rooms changed over time, as he became accustomed to them and the unfamiliar became familiar. Overall, the document is a reflection on the power of memory and how our experiences shape our perceptions of the world around us.
```

Using a JSON schema to get plant data

- Schema File:

```json
{
  "type": "object",
  "properties": {
    "plant-info": {
      "type": "object",
      "description": "Basic, factual details about the plant. Short answers.",
      "properties": {
        "common-name": {
          "type": "string"
        },
        "latin-name": {
          "type": "string"
        },
        "plant-description": {
          "type": "string"
        }
      },
      "required": [
        "common-name",
        "latin-name",
        "plant-description"
      ]
    },
    "planting-details": {
      "type": "object",
      "description": "Details on how to sow, care, and harvest the plant.",
      "properties": {
        "soil": {
          "type": "string"
        },
        "sun": {
          "type": "string"
        },
        "water": {
          "type": "string"
        },
        "sowing-method": {
          "type": "string"
        },
        "grow-days": {
          "type": "number"
        },
        "sowing-months": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "harvest-months": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "planting-instructions": {
          "type": "string"
        },
        "care-instructions": {
          "type": "string"
        },
        "harvest-instructions": {
          "type": "string"
        }
      },
      "required": [
        "soil",
        "sun",
        "water",
        "sowing-method",
        "grow-days",
        "sowing-months",
        "harvest-months",
        "planting-instructions",
        "care-instructionsharvest-instructions"
      ]
    }
  },
  "required": [
    "plant-info",
    "planting-details"
  ]
}
```

- Using the schema file. Note this currently requires a prompt

```bash
~ $: thyme -oa -model gpt4-0613 -c "As a gardener please give details about the following plant given your json schema" -a "broccoli"  -json ../test.txt
{
    "plant-info": {
        "common-name": "Broccoli",
        "latin-name": "Brassica oleracea",
        "plant-description": "Broccoli is a green vegetable that visually resembles a miniature tree. It is part of the Brassica oleracea species, the same family as cabbage, Brussels sprouts, and kale. Not only does it have a rich, savory taste, but it's also a nutrient powerhouse packed with vitamins, minerals, fiber, and antioxidants."
    },
    "planting-details": {
        "soil": "Rich, well-drained, slightly acidic soil",
        "sun": "Full sun",
        "water": "1 to 1.5 inches of water per week",
        "sowing-method": "Direct sowing or transplants",
        "grow-days": 80-100,
        "sowing-months": ["March","June"],
        "harvest-months": ["August","November"],
        "planting-instructions": "Sow seeds 0.5 inches deep and 3 inches apart. If using transplants, plant them 12-20 inches apart.",
        "care-instructions": "Keep the soil moist but not waterlogged. Provide consistent water supply, weeding, and protect from pests.",
        "harvest-instructions": "Harvest when the broccoli head is fully formed but before the tiny buds begin to open and yellow flowers emerge. Cut the stalk at an angle to prevent water from settling and rotting the remaining stalk."
    }
}
```

## External Resources
- [OpenAI Website](https://platform.openai.com/docs)
- [Kagi Universal Translator Website](https://kagi.com/summarizer/index.html)
