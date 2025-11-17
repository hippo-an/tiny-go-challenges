package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/gob"
	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/housework"
	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/json"
	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/protobuf"
)

var dataFile string
var t string

func init() {
	flag.StringVar(&dataFile, "file", "housework.db", "data file")
	flag.StringVar(&t, "type", "j|g|p", "JSON|Gob|Protocol Buffer serialize formatting")
	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			`Usage: %s [flags] [add chore, ... | complete #]
			add add comma-separated chores
			complete complete designated chore
			Flags:
			`,
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
}

func load() ([]*housework.Chore, error) {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return make([]*housework.Chore, 0), nil
	}

	df, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := df.Close(); err != nil {
			fmt.Printf("closing data file: %v", err)
		}
	}()

	switch t {
	case "j":
		return json.Load(df)
	case "g":
		return gob.Load(df)
	case "p":
		return protobuf.Load(df)
	default:
		return nil, fmt.Errorf("unsupported type: %s (use j, g, or p)", t)
	}
}

func flush(chores []*housework.Chore) error {
	df, err := os.Create(dataFile)
	if err != nil {
		return err
	}

	defer func() {
		if err := df.Close(); err != nil {
			fmt.Printf("closing data file: %v", err)
		}
	}()

	switch t {
	case "j":
		return json.Flust(df, chores)
	case "g":
		return gob.Flush(df, chores)
	case "p":
		return protobuf.Flush(df, chores)
	default:
		return fmt.Errorf("unsupported type: %s (use j, g, or p)", t)
	}
}

func list() error {
	chores, err := load()
	if err != nil {
		return err
	}

	if len(chores) == 0 {
		fmt.Println("You're all caught up!")
		return nil
	}
	for i, chore := range chores {
		status := " "
		if chore.Complete {
			status = "X"
		}
		fmt.Printf("[%s] %d: %s\n", status, i, chore.Description)
	}
	return nil
}

func add(s string) error {
	chores, err := load()
	if err != nil {
		return err
	}

	descriptions := strings.Split(s, ",")
	for _, desc := range descriptions {
		desc = strings.TrimSpace(desc)
		if desc != "" {
			chores = append(chores, &housework.Chore{
				Description: desc,
				Complete:    false,
			})
		}
	}

	return flush(chores)
}

func complete(arg string) error {
	chores, err := load()
	if err != nil {
		return err
	}

	index, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("invalid chore number: %s", arg)
	}

	if index < 0 || index >= len(chores) {
		return fmt.Errorf("chore number %d out of range (0-%d)", index, len(chores)-1)
	}

	chores[index-1].Complete = true

	return flush(chores)
}

func Run() error {
	flag.Parse()

	if t != "j" && t != "g" && t != "p" {
		return fmt.Errorf("invalid type: %s (use j for JSON, g for Gob, or p for Protobuf)", t)
	}

	args := flag.Args()

	if len(args) == 0 {
		return list()
	}

	switch strings.ToLower(args[0]) {
	case "add":
		if len(args) < 2 {
			return fmt.Errorf("add requires at least one chore description")
		}
		return add(strings.Join(args[1:], " "))
	case "complete":
		if len(args) < 2 {
			return fmt.Errorf("complete requires a chore number")
		}
		return complete(args[1])
	default:
		return fmt.Errorf("unknown command: %s (use 'add' or 'complete')", args[0])
	}

	err := list()
	if err != nil {
		log.Fatal(err)
	}

	return err
}
