BEGIN;

ALTER TABLE users
    ALTER COLUMN password_hash DROP NOT NULL;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_phone_verified BOOLEAN DEFAULT FALSE;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_phone_number_unique;

ALTER TABLE users
    ADD CONSTRAINT users_phone_number_unique UNIQUE (phone_number);

COMMIT;
