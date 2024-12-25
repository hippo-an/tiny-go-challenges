package utils

import (
	"bufio"
	"os"
)

func ReadFile(fileName string) ([]string, error) {
	var words []string
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}
