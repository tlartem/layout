BEGIN;

CREATE TYPE status AS ENUM (
    'pending',
    'active',
    'inactive',
    'banned'
    );

CREATE TABLE IF NOT EXISTS profile
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    name       TEXT,
    age        INT,
    status     status,
    verified   BOOLEAN
);

COMMIT;