package mailbox

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Mailbox struct {
	ID             uuid.UUID
	ownerID        uuid.UUID
	createdAt      time.Time
	lastListenedAt *time.Time
	opensAt        *time.Time
}

func NewMailbox(ownerId uuid.UUID, createdAt time.Time, opensAt *time.Time) (*Mailbox, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("mailbox: failed generating id: %w", err)
	}

	return &Mailbox{
		ID:             id,
		ownerID:        ownerId,
		createdAt:      createdAt,
		lastListenedAt: nil,
		opensAt:        opensAt,
	}, nil
}

func (m *Mailbox) SetOpeningTime(opensAt *time.Time) {
	m.opensAt = opensAt
}
