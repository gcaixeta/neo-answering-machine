package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gcaixeta/neo-answering-machine/mailbox"
	"github.com/google/uuid"
)

type MailboxHandler struct {
	repo mailbox.Repository
}

func (h *MailboxHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Retrieve the mailbox from the request
	var req struct {
		OwnerID uuid.UUID  `json:"owner_id"`
		OpensAt *time.Time `json:"opens_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	m, err := mailbox.NewMailbox(req.OwnerID, time.Now(), req.OpensAt)
	if err != nil {
		http.Error(w, "error creating mailbox", http.StatusInternalServerError)
		return
	}

	// Saves it to the database
	if err := h.repo.Save(r.Context(), m); err != nil {
		http.Error(w, "failed to save mailbox", http.StatusInternalServerError)
	}

	// Return the response with status code 201
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}
