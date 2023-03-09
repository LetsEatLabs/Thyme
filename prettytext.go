package main

import (
	"fmt"
	"regexp"
	"time"
)

// Removes any leading new lines before text in a string
func removeLeadingNewLines(s string) string {
	re := regexp.MustCompile(`^\s+`)
	return re.ReplaceAllString(s, "")
}

// Prints text so that it looks like a typewriter
func typeWriterPrint(s string) {

	re := regexp.MustCompile(`(?m)[^\S\r\n]{2,}`)
	newStr := re.ReplaceAllString(s, "")
	for _, c := range newStr {
		fmt.Printf("%c", c)
		time.Sleep(time.Millisecond * 25)
	}

	// One final space so we can separate lines printed in this fancy manner
	fmt.Printf(" ")
}
