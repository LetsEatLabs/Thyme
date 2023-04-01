package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

/////////////
/////////////

type ChatHistoryLine struct {
	Query  string
	Answer string
}

type ChatHistory struct {
	ChatHistoryLines []ChatHistoryLine
}

type SummaryHistory struct {
	Query  string
	Answer string
}

/////////////
/////////////

func viewHistoryQueries() {
	// Get the history file
	historyFiles := getHistoryFiles()

	fmt.Println(historyFiles)

}

/////////////

func getHistoryFiles() map[string][]string {

	historyFiles := make(map[string][]string)
	// Get the history file
	kagiDir := os.Getenv("THYME_QUERY_KAGI_LOGGING_DIR")
	openaiDir := os.Getenv("THYME_QUERY_LOGGING_DIR")

	historyFiles["kagi"] = getFilesInDir(kagiDir)
	historyFiles["openai"] = getFilesInDir(openaiDir)

	return historyFiles
}

/////////////

// Get a list of all files in a directory
func getFilesInDir(dir string) []string {
	filenames := []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			filenames = append(filenames, filepath.Join(dir, f.Name()))
		}
	}

	return filenames
}
