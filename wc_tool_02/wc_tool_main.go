package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func main() {
	// * cmd flag 설정
	byteFlag := flag.Bool("c", false, "Count bytes in a file")
	lineFlag := flag.Bool("l", false, "Count number of lines in a file")
	wordFlag := flag.Bool("w", false, "Count number of words in a file")
	flag.Parse()

	// content 받는 방법 2가지
	// 1. stdin -> pipe
	// 2. read file -> args
	content, filePath := parseInput()

	if *byteFlag {
		fmt.Printf("%8d %s\n", len(content), filePath)
	} else if *lineFlag {
		fmt.Printf("%8d %s\n", countLineInString(string(content)), filePath)
	} else if *wordFlag {
		fmt.Printf("%8d %s\n", countWordsInString(string(content)), filePath)
	} else {
		fmt.Printf("%8d %8d %8d %s\n", countLineInString(string(content)), countWordsInString(string(content)), len(content), filePath)
	}
}

func parseInput() ([]byte, string) {
	filePath := ""
	var content []byte
	if flag.NArg() == 0 {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			exit("Error reading input:", err)
		}
		content = bytes
	} else {
		filePath = flag.Arg(0)
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			exit("Error reading file:", err)
		}
		content = bytes
	}

	return content, filePath
}

func countLineInString(content string) int {
	return len(strings.Split(content, "\n"))
}

func countWordsInString(s string) int {
	var wordCount int
	inWord := false

	for _, char := range s {
		if unicode.IsSpace(char) {
			inWord = false
		} else if !inWord {
			wordCount++
			inWord = true
		}
	}

	return wordCount
}

func exit(messages ...any) {
	fmt.Println(messages)
	os.Exit(1)
}
