package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/lipgloss"
)

/////////////////
// Lipgloss Styles

var (
	// Lipgloss style for the query waiting spinner
	spinnerText = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
)

/////////////////

// Removes any leading new lines before text in a string
func removeLeadingNewLines(s string) string {
	re := regexp.MustCompile(`^\s+`)
	return re.ReplaceAllString(s, "")
}

/////////////////

// Prints text so that it looks like a typewriter
// GPT-4 Helped fixed the colors printing slowly
func typeWriterPrint(s string, space bool) {
	re := regexp.MustCompile(`(?m:!^)[^\S\r\n\t]{2,}`)
	newStr := re.ReplaceAllString(s, "")

	// Use a regular expression to match the terminal's escape sequences
	// (\033 represents the ESC character in octal notation)
	ansiEscapeSeq := regexp.MustCompile(`\033\[[0-9;]*m`)

	// Keep track of whether we have an ongoing escape sequence
	// This is so that text that is supposed to be in color
	// Prints at the same speed at regular text
	inEscapeSequence := false
	escSeq := ""

	for _, c := range newStr {
		str := string(c)
		if inEscapeSequence {
			escSeq += str
			if ansiEscapeSeq.MatchString(escSeq) {
				inEscapeSequence = false
				fmt.Print(escSeq)
			}
		} else {
			if str == "\033" {
				inEscapeSequence = true
				escSeq = str
			} else {
				fmt.Printf("%c", c)
				time.Sleep(20 * time.Millisecond)
			}
		}
	}

	if space {
		// One final space so we can separate lines printed in this fancy manner
		fmt.Printf(" ")
	} else {
		fmt.Printf("\n")
	}
}

/////////////////

// Function that is a spinner that last until a query is done
func spinner(spinningComplete chan bool) {
	for {

		for _, r := range `▁▂▃▄▅▆▇█▇▆▅▄▃▂▁` {
			prettyPrintSpinner(fmt.Sprintf("\r%c Querying...", r))
			time.Sleep(time.Millisecond * 100)
		}

		select {
		case value := <-spinningComplete:
			if value == true {
				fmt.Printf("\r                    \r")
				return
			}

		default:
			continue
		}
	}
}

/////////////////

// Display the available prompts to the user
func listAvailablePrompts() {

	prompts := initPrompts()

	fmt.Printf("Available prompts:\n\n")
	for _, prompt := range prompts {
		fmt.Printf("- %s: %s\n", prompt.Name, prompt.Description)
	}
}

/////////////////

// Pretty print the text green
func prettyPrintSpinner(s string) {
	fmt.Printf("%s", spinnerText.Render(s))
}

// Pretty print the text green
func prettyPrintChatArrow(s string) {
	fmt.Printf("%s", spinnerText.Render(s))
}

/////////////////

// Use the Chroma library to guess the syntax of the string and format it to print
// To the terminal
// This function was co-authored with GPT-4
func prettyPrintCode(s string, language string) string {
	var lexer chroma.Lexer

	if language != "" {
		lexer = lexers.Get(language)
	}

	// If we did not pass a language, then please guess
	if language == "" {
		detectProgrammingLanguage(s)
	}

	if lexer == nil {
		lexer = lexers.Analyse(s)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	// If we can guess the language, print it with syntax highlighting
	//lexer = chroma.Coalesce(lexer)

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// Creating a bytes.Buffer to hold the formatted string
	var buf bytes.Buffer

	// Formatting and printing the syntax-highlighted code
	iterator, _ := lexer.Tokenise(nil, s)
	err := formatter.Format(&buf, style, iterator)
	if err != nil {
		fmt.Println("Error formatting code:", err)
	}

	// Returning the string representation of the buffer
	return buf.String()

}

/////////////////

// Format just the codeblocks in Markdown
// Function written by GPT-4
func formatCodeBlocksInMarkdown(s string, language string) string {
	// Regular expression to find code blocks: ``` followed by optional language,
	// then anything until another ```
	codeBlockRegex := regexp.MustCompile("(?s)(```)(.*?)(```)")

	result := codeBlockRegex.ReplaceAllStringFunc(s, func(match string) string {
		submatches := codeBlockRegex.FindStringSubmatch(match)
		code := submatches[2] // Extract the code portion from the match

		formattedCode := prettyPrintCode(code, language)

		// Replacing the original code with the formatted code while preserving
		// the original triple backticks and optional language identifier
		return codeBlockRegex.ReplaceAllString(match, "$1"+formattedCode+"$1")
	})

	return result
}

/////////////////

// Try to detect the language from the response
// Made with help from GPT-4

func detectProgrammingLanguage(text string) string {
	keywords := map[string][]string{
		"python":     {"def", "import", "print", "class"},
		"java":       {"public", "static", "void", "class"},
		"go":         {"func", "import", "package", "type"},
		"javascript": {"function", "var", "let", "const"},
		"ruby":       {"def", "class", "require", "end"},
		"bash":       {"#!/bin/bash", "echo", "if", "fi"},
		"cpp":        {"#include", "iostream", "using", "namespace"},
		"c":          {"#include", "stdio.h", "int", "main"},
		"csharp":     {"using", "namespace", "class", "public"},
		"rust":       {"fn", "mut", "let", "use"},
		"scala":      {"object", "def", "val", "var"},
		"ada":        {"procedure", "is", "begin", "end"},
		"gdscript":   {"extends", "func", "var", "pass"},
		"perl":       {"#!/usr/bin/perl", "use", "strict", "print"},
	}

	golangLibraries := []string{"fmt", "os", "math", "http", "time", "json"}
	javascriptObjects := []string{"document", "console", "window", "Array", "Date"}
	pythonObjects := []string{"sys", "os", "math", "random", "datetime", "__name__"}
	gdscriptObjects := []string{"Node", "Spatial", "KinematicBody", "func", "var"}

	lowPriority := ""
	for lang, words := range keywords {
		for _, word := range words {
			regex := regexp.MustCompile(`\b` + word + `\b`)
			if regex.MatchString(text) {

				// Check if it is golang first
				if lang == "go" {
					for _, lib := range golangLibraries {
						if strings.Contains(text, "\""+lib+"\"") {
							return lang
						}
					}
				}

				// Javascript
				if lang == "javascript" {
					for _, obj := range javascriptObjects {
						if strings.Contains(text, obj) {
							return lang
						}
					}
				}

				// GDScript
				if lang == "gdscript" {
					for _, obj := range gdscriptObjects {
						if strings.Contains(text, obj) {
							return lang
						}
					}
				}

				// Python
				if lang == "python" {
					for _, obj := range pythonObjects {
						if strings.Contains(text, obj) {
							return lang
						}
					}
				}

				return lang
			}
		}
	}
	return lowPriority
}
