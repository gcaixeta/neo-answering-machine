// Package api declares a web api where the user can call the main functionalities from Neo
package api

import (
	"net/http"

	"github.com/gcaixeta/neo-answering-machine/mailbox"
	"github.com/gcaixeta/neo-answering-machine/tape"
)

func NewRouter(mailboxes mailbox.Repository, tapes tape.Repository) *http.ServeMux {
	mux := http.NewServeMux()

	mh := &MailboxHandler{repo: mailboxes}
	th := &TapeHandler{repo: tapes}

	mux.HandleFunc("POST /mailbox", mh.Create)

	mux.HandleFunc("GET /mailbox/{id}", mh.GetByID)

	mux.HandleFunc("POST /tape/upload", th.UploadNewTape)

	return mux
}
