CREATE TABLE IF NOT EXISTS mailboxes (
    id               UUID PRIMARY KEY,
    owner_id         UUID NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL,
    last_listened_at TIMESTAMPTZ,
    opens_at         TIMESTAMPTZ
);
