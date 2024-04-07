package main

import (
	"io"
	"os"
)

type SlugReader interface {
	Read(slug string) (string, error)
}

type MdFileReader struct{}

func (fsr MdFileReader) Read(slug string) (string, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
