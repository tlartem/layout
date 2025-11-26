BEGIN;

ALTER TABLE profile
    ADD COLUMN contacts JSONB; -- email, phone

COMMIT;