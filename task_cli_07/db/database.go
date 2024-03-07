package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Task struct {
	Id   int64
	Task string
}

func InitDB(file string) error {
	var err error
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	createTableQuery := `
		create table IF NOT EXISTS task ( 
		id integer PRIMARY KEY autoincrement,
		task text,
		UNIQUE (id)
		)
	`

	if _, err = db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}
func DeleteTask(id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("delete from task where id = ?", id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func AllTasks() ([]Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT * from task")
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func CreateTask(task string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO task VALUES (NULL, ?)", task)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return id, nil
}
