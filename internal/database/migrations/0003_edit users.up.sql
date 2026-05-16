BEGIN ;
ALTER TABLE users
    ALTER COLUMN phone_number SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT users_phone_number_unique UNIQUE(phone_number);

COMMIT ;