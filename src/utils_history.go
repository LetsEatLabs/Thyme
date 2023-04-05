package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/////////////
/////////////

type ChatHistoryLine struct {
	Query  string `json:"query"`
	Answer string `json:"answer"`
}

type QueryHistory struct {
	Query  string `json:"query"`
	Answer string `json:"answer"`
}

type ChatHistory struct {
	ChatHistoryLines []ChatHistoryLine
}

type SummaryHistory struct {
	Query  string `json:"query"`
	Answer string `json:"answer"`
}

/////////////
/////////////

func viewHistoryQueries(historyFlag string) {
	// Get the history file
	historyFiles := getHistoryFiles()

	// Get the font styles
	styles := getFontStyles()

	// Check if we want to view a file first
	if doesFileExist(historyFlag) {

		if strings.Contains(historyFlag, "query") {
			queryHistory := loadQueryHistoryFile(historyFlag)
			fmt.Println(styles.historyTitle.Render("Query: "))
			fmt.Println(styles.historyTitle.Render("----------"))
			fmt.Println(queryHistory.Query)
			fmt.Println()

			fmt.Println(styles.historyTitle.Render("Answer: "))
			fmt.Println(styles.historyTitle.Render("----------"))
			content := formatCodeBlocksInMarkdown(queryHistory.Answer, "")

			fmt.Println(content)
		}

		if strings.Contains(historyFlag, "chat") {
			queryHistory := loadChatHistoryFile(historyFlag)

			chl := queryHistory.ChatHistoryLines
			for i := range chl {
				fmt.Println(styles.historyTitle.Render("Query: "))
				fmt.Println(styles.historyTitle.Render("----------"))
				fmt.Println(chl[i].Query)
				fmt.Println()

				fmt.Println(styles.historyTitle.Render("Answer: "))
				fmt.Println(styles.historyTitle.Render("----------"))
				content := formatCodeBlocksInMarkdown(chl[i].Answer, "")

				fmt.Println(content)
				fmt.Println()
			}

		}

		if strings.Contains(historyFlag, "summary") {
			queryHistory := loadSummaryHistoryFile(historyFlag)
			fmt.Println(styles.historyTitle.Render("Source: "))
			fmt.Println(styles.historyTitle.Render("----------"))
			fmt.Println(queryHistory.Query)
			fmt.Println()

			fmt.Println(styles.historyTitle.Render("Summary: "))
			fmt.Println(styles.historyTitle.Render("----------"))
			content := formatCodeBlocksInMarkdown(queryHistory.Answer, "")

			fmt.Println(content)
		}

		os.Exit(0)
	}

	if historyFlag == "query" || historyFlag == "all" {
		fmt.Println()
		fmt.Println(styles.historyTitle.Render("Queries"))
		fmt.Println(styles.historyTitle.Render("----------"))

		// Load the history files
		for i := range historyFiles["openai"] {
			fname := historyFiles["openai"][i]
			if strings.Contains(fname, "query") {
				queryHistory := loadQueryHistoryFile(fname)
				if len(queryHistory.Query) > 75 {
					queryHistory.Query = queryHistory.Query[:75]
				}
				fmt.Println(styles.historyInfo.Render("File: ") + styles.historyText.Render(fname))
				fmt.Println(styles.historyInfo.Render("Asked: ") + styles.historyText.Render(queryHistory.Query) + "\n")
			}
		}
	}

	if historyFlag == "summary" || historyFlag == "all" {

		fmt.Println()
		fmt.Println(styles.historyTitle.Render("Summaries"))
		fmt.Println(styles.historyTitle.Render("----------"))

		// Load the summary files
		for i := range historyFiles["kagi"] {
			fname := historyFiles["kagi"][i]
			if strings.Contains(fname, "summary") {
				summaryHistory := loadSummaryHistoryFile(fname)
				if len(summaryHistory.Query) > 75 {
					summaryHistory.Query = summaryHistory.Query[:75]
				}
				fmt.Println(styles.historyInfo.Render("File: ") + styles.historyText.Render(fname))
				fmt.Println(styles.historyInfo.Render("Summarized: ") + styles.historyText.Render(summaryHistory.Query) + "\n")
			}
		}
	}

	if historyFlag == "chat" || historyFlag == "all" {

		fmt.Println()
		fmt.Println(styles.historyTitle.Render("Chats"))
		fmt.Println(styles.historyTitle.Render("----------"))

		// Load the chat files
		for i := range historyFiles["openai"] {
			fname := historyFiles["openai"][i]
			if strings.Contains(fname, "chat") {

				chat := loadChatHistoryFile(fname)
				targetChat := chat.ChatHistoryLines[0]
				if len(targetChat.Query) > 75 {
					targetChat.Query = targetChat.Query[:75]
				}

				fmt.Println(styles.historyInfo.Render("File: ") + styles.historyText.Render(fname))
				fmt.Println(styles.historyInfo.Render("Starter: ") + styles.historyText.Render(targetChat.Query) + "\n")
			}
		}
	}

	//fmt.Println(historyFiles)

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

/////////////

// Load a query history file, and return a QueryHistory object
func loadQueryHistoryFile(filename string) QueryHistory {
	filestr := readFileToString(filename)
	queryHistory := QueryHistory{}
	json.Unmarshal([]byte(filestr), &queryHistory)
	return queryHistory
}

/////////////

// Load a summary history file, and return a SummaryHistory object
func loadSummaryHistoryFile(filename string) SummaryHistory {
	filestr := readFileToString(filename)
	summaryHistory := SummaryHistory{}
	json.Unmarshal([]byte(filestr), &summaryHistory)
	return summaryHistory
}

/////////////

// Load all ChatHitoryLines in a chat history file, and return a ChatHistory object
// Cowritten by GPT-4
func loadChatHistoryFile(filename string) ChatHistory {
	chatHistory := ChatHistory{ChatHistoryLines: []ChatHistoryLine{}}

	// Open the JSONL file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file line-by-line and append each line into a string slice
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Iterate through the slice, unmarshal each JSON object string and process it
	for _, lineJson := range lines {
		var obj ChatHistoryLine
		err := json.Unmarshal([]byte(lineJson), &obj)
		if err != nil {
			fmt.Printf("failed to unmarshal json: %s\n", err)
			os.Exit(1)
		}

		chatHistory.ChatHistoryLines = append(chatHistory.ChatHistoryLines, obj)
	}

	return chatHistory
}

/////////////

// Checks if a string is a file that exists on the system
func doesFileExist(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	} else {
		return !fileInfo.IsDir() // We only want files, not directories
	}
}

/////////////
