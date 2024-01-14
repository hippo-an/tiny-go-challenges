package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()

	file, err := readFile(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file:%s.\n", *csvFileName))
	}
	defer file.Close()

	lines, err := readCSV(file)
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseToProblem(lines)
	correctCount := uint64(0)
	scanner := bufio.NewScanner(os.Stdin)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	answerChan := make(chan string)

	defer func() {
		timer.Stop()
		close(answerChan)
	}()

problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		go func() {
			scanner.Scan()
			answer := strings.TrimSpace(scanner.Text())
			answerChan <- answer
		}()

		select {
		case <-timer.C:
			break problemLoop
		case answer := <-answerChan:
			if answer == p.answer {
				correctCount++
			}
		}
	}

	fmt.Printf("\nYou scored %d out of %d.\n", correctCount, len(problems))
}

func parseToProblem(lines [][]string) []problem {
	problems := make([]problem, len(lines))

	for i, r := range lines {
		problem := problem{
			question: strings.TrimSpace(r[0]),
			answer:   strings.TrimSpace(r[1]),
		}
		problems[i] = problem
	}

	return problems
}

func readCSV(file io.Reader) ([][]string, error) {
	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}

	return rows, nil
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func readFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
