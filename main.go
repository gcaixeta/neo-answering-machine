package main

import (
	"log"
	"net/http"

	"github.com/gcaixeta/neo-answering-machine/internal/api"
	"github.com/gcaixeta/neo-answering-machine/internal/repository/postgres"
	_ "github.com/lib/pq"
)

func main() {
	const dsn = "postgres://neo:neo@localhost:5439/neo?sslmode=disable"

	db, err := postgres.NewDB(dsn)
	if err != nil {
		log.Fatalf("error opening db connection: %v", err)
	}
	defer db.Close()

	mailboxRepo := postgres.NewMailboxRepository(db)
	tapeRepo := postgres.NewTapeRepository(db)
	mux := api.NewRouter(mailboxRepo, tapeRepo)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
