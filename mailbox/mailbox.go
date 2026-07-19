// Package mailbox implements the mailbox used to store the tapes sent for a unique user.
//
// A Mailbox belongs to a single owner and holds tapes recorded for them. Each
// Mailbox has an opensAt time that determines when its tapes become available:
// before that time, the owner's latest unlistened tapes stay locked. A nil
// opensAt means the mailbox is open with no delay.
//
// LastListenedAt is nil until the owner listens to a tape for the first time.package mailbox
package mailbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("mailbox: not found")

type Mailbox struct {
	id             uuid.UUID
	ownerID        uuid.UUID
	createdAt      time.Time
	lastListenedAt *time.Time
	opensAt        *time.Time
}

type Repository interface {
	Save(ctx context.Context, mb *Mailbox) error

	FindByID(ctx context.Context, id uuid.UUID) (*Mailbox, error)
}

func NewMailbox(ownerID uuid.UUID, createdAt time.Time, opensAt *time.Time) (*Mailbox, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("mailbox: failed generating id: %w", err)
	}

	return &Mailbox{
		id:             id,
		ownerID:        ownerID,
		createdAt:      createdAt,
		lastListenedAt: nil,
		opensAt:        opensAt,
	}, nil
}

func (m *Mailbox) SetOpeningTime(opensAt *time.Time) {
	m.opensAt = opensAt
}

func (m *Mailbox) ID() uuid.UUID {
	return m.id
}

func (m *Mailbox) OwnerID() uuid.UUID {
	return m.ownerID
}

func (m *Mailbox) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Mailbox) LastListenedAt() *time.Time {
	return m.lastListenedAt
}

func (m *Mailbox) OpensAt() *time.Time {
	return m.opensAt
}

func Reconstruct(id, ownerID uuid.UUID, createdAt time.Time, lastListenedAt, opensAt *time.Time) *Mailbox {
	return &Mailbox{
		id:             id,
		ownerID:        ownerID,
		createdAt:      createdAt,
		lastListenedAt: lastListenedAt,
		opensAt:        opensAt,
	}
}
