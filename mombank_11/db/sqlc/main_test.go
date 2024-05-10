package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable"
)

var testQueries *Queries

// package 에서 하나의 함수에서만 사용할 수 있는 test entry 를 설정
// main 을 수행하고 test 를 수행한다.
func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = conn.Ping()

	if err != nil {
		log.Fatal("connection error:", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
