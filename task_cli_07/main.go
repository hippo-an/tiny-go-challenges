package main

import (
	"log"
	"path/filepath"

	"github.com/dev-hippo-an/tiny-go-challenges/task_cli_07/cmd"
	"github.com/dev-hippo-an/tiny-go-challenges/task_cli_07/db"
	"github.com/mitchellh/go-homedir"
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal("error while get homedir")
	}
	dbPath := filepath.Join(home, "task_test.db")
	// log.Println(home)
	// dbPath := "file::memory:?mode=memory&cache=shared"

	err = db.InitDB(dbPath)
	if err != nil {
		log.Fatal("sqlite3 db is not connected.", err)
	}

	cmd.Execute()
}
