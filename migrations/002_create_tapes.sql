CREATE TABLE IF NOT EXISTS tapes (
    id          UUID PRIMARY KEY,
    mailbox_id  UUID NOT NULL REFERENCES mailboxes(id),
    recorded_by UUID NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL,
    played      BOOLEAN NOT NULL DEFAULT false,
    playedAt    TIMESTAMPTZ
);
