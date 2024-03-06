package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Task struct {
	Key   int
	Value string
}

func InitDB(file string) error {
	var err error
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

	db.BeginTx(context.Background(), &sql.TxOptions{})

	createTableQuery := `
		create table IF NOT EXISTS task ( 
		id integer PRIMARY KEY autoincrement,
		userId text,
		password text,
		UNIQUE (id, userId)
		)
	`

	if _, err = db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}
