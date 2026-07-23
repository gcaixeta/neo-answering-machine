// Package postgres contains implementations for neo repos.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gcaixeta/neo-answering-machine/mailbox"
	"github.com/google/uuid"
)

type MailboxRepository struct {
	db *sql.DB
}

type scanner interface {
	Scan(dest ...any) error
}

func NewMailboxRepository(db *sql.DB) *MailboxRepository {
	return &MailboxRepository{db: db}
}

var _ mailbox.Repository = (*MailboxRepository)(nil)

func (r *MailboxRepository) Save(ctx context.Context, mb *mailbox.Mailbox) error {
	_, err := r.db.ExecContext(ctx, `
			INSERT INTO mailboxes (id, owner_id, created_at, last_listened_at, opens_at)
			VALUES ($1, $2, $3, $4, $5)
		`, mb.ID(), mb.OwnerID(), mb.CreatedAt(), mb.LastListenedAt(), mb.OpensAt())
	if err != nil {
		return fmt.Errorf("mailbox: save: %w", err)
	}

	return nil
}

func (r *MailboxRepository) FindByID(ctx context.Context, id uuid.UUID) (*mailbox.Mailbox, error) {
	row := r.db.QueryRow(`
			SELECT id, owner_id, created_at, last_listened_at, opens_at FROM mailboxes WHERE id = $1
		`, id)

	m, err := scanMailbox(row)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, mailbox.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("mailbox: FindByID: %w", err)
	}

	return m, nil
}

func scanMailbox(s scanner) (*mailbox.Mailbox, error) {
	var (
		id, ownerID             uuid.UUID
		createdAt               time.Time
		LastListenedAt, opensAt *time.Time
	)
	if err := s.Scan(&id, &ownerID, &createdAt, &LastListenedAt, &opensAt); err != nil {
		return nil, err
	}
	return mailbox.Reconstruct(id, ownerID, createdAt, LastListenedAt, opensAt), nil
}
