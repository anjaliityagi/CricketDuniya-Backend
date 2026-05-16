BEGIN;


TRUNCATE TABLE users CASCADE;

ALTER TABLE users
    ADD COLUMN phone_number VARCHAR(15);

ALTER TABLE users
    DROP COLUMN email;



CREATE TABLE IF NOT EXISTS user_session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   user_id UUID NOT NULL REFERENCES users(id),
     created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  archived_at TIMESTAMP WITH TIME ZONE
);

COMMIT;