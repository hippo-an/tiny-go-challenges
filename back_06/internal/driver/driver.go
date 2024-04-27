package driver

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

func ConnectSQL(driverName, dataSourceName string) (*DB, error) {
	d, err := NewDatabase(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = d

	err = testDb()
	if err != nil {
		panic(err)
	}

	return dbConn, nil
}

func testDb() error {
	if err := dbConn.SQL.Ping(); err != nil {
		return err
	}

	return nil
}

func NewDatabase(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
