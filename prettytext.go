package main

import (
	"fmt"
	"regexp"
	"time"

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
func typeWriterPrint(s string) {

	re := regexp.MustCompile(`(?m:!^)[^\S\r\n\t]{2,}`)
	newStr := re.ReplaceAllString(s, "")
	for _, c := range newStr {
		fmt.Printf("%c", c)
		time.Sleep(time.Millisecond * 20)
	}

	// One final space so we can separate lines printed in this fancy manner
	fmt.Printf(" ")
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

/////////////////
