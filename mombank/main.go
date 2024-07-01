package main

import (
	"database/sql"
	"log"

	"github.com/hippo-an/tiny-go-challenges/mombank/api"
	db "github.com/hippo-an/tiny-go-challenges/mombank/db/sqlc"
	"github.com/hippo-an/tiny-go-challenges/mombank/util"

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

	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
