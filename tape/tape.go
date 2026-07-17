package tape

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Tape struct {
	id         uuid.UUID
	mailboxId  uuid.UUID
	recordedBy uuid.UUID //references user id
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
		id:         id,
		mailboxId:  mailboxId,
		recordedBy: recordedBy,
		recordedAt: recordedAt,
		played:     false,
		playedAt:   nil,
	}, nil
}
