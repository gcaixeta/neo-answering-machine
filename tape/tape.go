// Package tape implements the tapes stored inside a user's mailbox.
//
// A Tape represents a single voice message recorded by one user for another
// mailbox's owner. Each Tape tracks who recorded it, when, and whether the
// owner has played it yet. PlayedAt is nil until the tape is first played.
package tape

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Tape struct {
	ID         uuid.UUID
	mailboxID  uuid.UUID
	recordedBy uuid.UUID // references user id
	recordedAt time.Time
	played     bool
	playedAt   *time.Time
}

func NewTape(mailboxId, recordedBy uuid.UUID, recordedAt time.Time) (*Tape, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("tape: failed to generate id: %w", err)
	}

	return &Tape{
		ID:         id,
		mailboxID:  mailboxId,
		recordedBy: recordedBy,
		recordedAt: recordedAt,
		played:     false,
		playedAt:   nil,
	}, nil
}
