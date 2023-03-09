# Thyme

A CLI helper for interacting with the ChatGPT API. Capable of sending files using built-in prompts, examples (coming soon), and prompt-chains (coming soon). Also capable of simple direct queries.

## Installation

```bash
~ $: export OPENAI_API_KEY="<your_key_here>"
~ $: go build -o thyme *.go
~ $: ./thyme --help

Usage: thyme <flags> <input file>

Flags:
    -p <prompt>     The prompt to use for the GPT request
    -q <question>   Ask a question and get a response (do not use with -p)
    -h (--help)     Display this help message
```

## Examples

```bash
~ $: thyme -q "Can you explain how a RICO case works in the
 voice of a baseball announcer"
Welcome to today's game folks, and we've got a real interesting one for you. This matchup features a Racketeer Influenced and Corrupt Organizations Act case, or RICO for short. 

Now, to break it down for you - this is when the government charges a group of individuals or an organization with committing a pattern of crimes. Think of it like a baseball team that's been caught cheating by stealing signs or juicing up their players with PEDs. 

The goal of the government is to dismantle this criminal enterprise by taking down the key players and seizing their assets. Just like a manager would make strategic moves to bring in the right players to win the game, the government's attorneys and investigators have to work together to prove their case against the accused.

Players in a RICO case can face hefty fines and prison time, just like a player caught using steroids would face suspension and tarnished reputation. 

So, sit back folks and get ready for a real nail biter of a case. It's not often we see a RICO case on the field, but when we do, it's sure to be a game-changer.
```

```bash
~ $: thyme -p listify test.txt
1. Go to the store and buy eggs and milk.
2. Stop by the gasoline store and buy tires.
3. Purchase a whole salmon and eat it.
```
