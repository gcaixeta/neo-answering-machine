package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gcaixeta/neo-answering-machine/tape"
	"github.com/google/uuid"
)

type TapeRepository struct {
	db *sql.DB
}

func NewTapeRepository(db *sql.DB) *TapeRepository {
	return &TapeRepository{db: db}
}

var _ tape.Repository = (*TapeRepository)(nil)

func (r *TapeRepository) Save(ctx context.Context, t *tape.Tape) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tapes (id, mailbox_id, recorded_by, recorded_at, played, playedAt)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, t.ID, t.MailboxID(), t.RecordedBy(), t.RecordedAt(), t.Played(), t.PlayedAt())
	if err != nil {
		return fmt.Errorf("tape: save: %w", err)
	}

	return nil
}

func (r *TapeRepository) FindByID(ctx context.Context, id uuid.UUID) (*tape.Tape, error) {
	row := r.db.QueryRow(`
		SELECT id, mailbox_id, recorded_by, recorded_at, played, playedAt
		FROM tapes
		WHERE id = $1
	`, id)

	t, err := scanTape(row)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, tape.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("tape: FindByID: %w", err)
	}

	return t, nil
}

func (r *TapeRepository) ListByMailboxID(ctx context.Context, mailboxID uuid.UUID) ([]*tape.Tape, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, mailbox_id, recorded_by, recorded_at, played, playedAt
		FROM tapes
		WHERE mailbox_id = $1
	`, mailboxID)
	if err != nil {
		return nil, fmt.Errorf("tape: ListByMailboxID: %w", err)
	}
	defer rows.Close()

	var tapes []*tape.Tape
	for rows.Next() {
		t, err := scanTape(rows)
		if err != nil {
			return nil, fmt.Errorf("tape: ListByMailboxID: %w", err)
		}
		tapes = append(tapes, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("tape: ListByMailboxID: %w", err)
	}

	return tapes, nil
}

func scanTape(s scanner) (*tape.Tape, error) {
	var (
		id, mailboxID, recordedBy uuid.UUID
		recordedAt                time.Time
		played                    bool
		playedAt                  *time.Time
	)

	if err := s.Scan(&id, &mailboxID, &recordedBy, &recordedAt, &played, &playedAt); err != nil {
		return nil, err
	}

	return tape.Reconstruct(id, mailboxID, recordedBy, recordedAt, played, playedAt), nil
}
