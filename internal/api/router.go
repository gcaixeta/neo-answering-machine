// Package api declares a web api where the user can call the main functionalities from Neo
package api

import (
	"net/http"

	"github.com/gcaixeta/neo-answering-machine/mailbox"
)

func NewRouter(mailboxes mailbox.Repository) *http.ServeMux {
	mux := http.NewServeMux()

	h := &MailboxHandler{repo: mailboxes}
	mux.HandleFunc("POST /mailbox", h.Create)

	mux.HandleFunc("GET /mailbox/{id}", h.GetByID)

	return mux
}
