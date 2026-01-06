package main

import (
	"log"

	"ontopsolutions.net/gasperlf/social/internal/db"
	"ontopsolutions.net/gasperlf/social/internal/env"
	"ontopsolutions.net/gasperlf/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:pass@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 10, 10, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store)
}
