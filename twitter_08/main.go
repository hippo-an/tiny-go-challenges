package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	var keys struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}

	f, err := os.Open(".keys.json")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&keys)
	if err != nil {
		return
	}

	log.Println(keys)

}
