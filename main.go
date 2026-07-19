package neoansweringmachine

import (
	"fmt"

	"github.com/gcaixeta/neo-answering-machine/internal/api"
	"github.com/gcaixeta/neo-answering-machine/internal/repository/postgres"
)

func main() {
	const dsn = "blablabla"

	db, err := postgres.NewDB(dsn)
	if err != nil {
		fmt.Printf("Error trying to open connection with db: %w", err)
	}

	mailboxRepo := postgres.NewMailboxRepository(db)
	api.NewRouter(mailboxRepo)
}
