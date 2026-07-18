// Package tape implements the tapes stored inside a user's mailbox.
//
// A Tape represents a single voice message recorded by one user for another
// mailbox's owner. Each Tape tracks who recorded it, when, and whether the
// owner has played it yet. PlayedAt is nil until the tape is first played.
package tape

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("tape: not found")

type Tape struct {
	ID         uuid.UUID
	mailboxID  uuid.UUID
	recordedBy uuid.UUID // references user id
	recordedAt time.Time
	played     bool
	playedAt   *time.Time
}

type Repository interface {
	Save(ctx context.Context, tape *Tape) error

	FindByID(ctx context.Context, id uuid.UUID) (*Tape, error)

	ListByMailboxID(ctx context.Context, mailboxID uuid.UUID) ([]*Tape, error)
}

func NewTape(mailboxID, recordedBy uuid.UUID, recordedAt time.Time) (*Tape, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("tape: failed to generate id: %w", err)
	}

	return &Tape{
		ID:         id,
		mailboxID:  mailboxID,
		recordedBy: recordedBy,
		recordedAt: recordedAt,
		played:     false,
		playedAt:   nil,
	}, nil
}

func (t *Tape) MarkPlayed(playedAt time.Time) {
	t.playedAt = &playedAt
	t.played = true
}

func (t *Tape) MailboxID() uuid.UUID {
	return t.mailboxID
}

func (t *Tape) RecordedBy() uuid.UUID {
	return t.recordedBy
}

func (t *Tape) RecordedAt() time.Time {
	return t.recordedAt
}

func (t *Tape) Played() bool {
	return t.played
}

func (t *Tape) PlayedAt() *time.Time {
	return t.playedAt
}
