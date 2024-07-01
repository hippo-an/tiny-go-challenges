package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/hippo-an/tiny-go-challenges/mombank/util"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// package 에서 하나의 함수에서만 사용할 수 있는 test entry 를 설정
// main 을 수행하고 test 를 수행한다.
func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal(err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = testDB.Ping()

	if err != nil {
		log.Fatal("connection error:", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
