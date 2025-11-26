BEGIN;

CREATE TABLE IF NOT EXISTS property
(
    profile_id UUID REFERENCES profile (id) ON DELETE CASCADE,
    tags       TEXT[]
);

COMMIT;