package api

import (
	"time"

	"github.com/gcaixeta/neo-answering-machine/mailbox"
	"github.com/google/uuid"
)

type mailboxResponse struct {
	ID             uuid.UUID  `json:"id"`
	OwnerID        uuid.UUID  `json:"owner_id"`
	CreatedAt      time.Time  `json:"created_at"`
	LastListenedAt *time.Time `json:"last_listened_at"`
	OpensAt        *time.Time `json:"opens_at"`
}

func newMailboxResponse(m *mailbox.Mailbox) mailboxResponse {
	return mailboxResponse{
		ID:             m.ID(),
		OwnerID:        m.OwnerID(),
		CreatedAt:      m.CreatedAt(),
		LastListenedAt: m.LastListenedAt(),
		OpensAt:        m.OpensAt(),
	}
}
