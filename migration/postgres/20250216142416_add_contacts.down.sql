BEGIN;

ALTER TABLE profile
    DROP COLUMN contacts;

COMMIT;