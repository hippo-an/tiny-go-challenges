package main

import (
	"os"

	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
