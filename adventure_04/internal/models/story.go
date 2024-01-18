package models

import (
	"encoding/json"
	"io"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title,omitempty"`
	Paragraphs []string `json:"story,omitempty"`
	Options    []Option `json:"options,omitempty"`
}

type Option struct {
	Text    string `json:"text,omitempty"`
	Chapter string `json:"arc,omitempty"`
}

func JsonStory(r io.Reader) (Story, error) {
	decoder := json.NewDecoder(r)
	var story Story
	err := decoder.Decode(&story)
	if err != nil {
		return nil, err
	}
	return story, nil
}
