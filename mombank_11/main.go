package main

import (
	"database/sql"
	"github.com/dev-hippo-an/tiny-go-challenges/mombank_11/api"
	db "github.com/dev-hippo-an/tiny-go-challenges/mombank_11/db/sqlc"
	"github.com/dev-hippo-an/tiny-go-challenges/mombank_11/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = conn.Ping()

	if err != nil {
		log.Fatal("connection error:", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
