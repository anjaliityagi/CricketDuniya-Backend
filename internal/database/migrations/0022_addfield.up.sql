BEGIN ;


ALTER TABLE innings
    ADD COLUMN target_runs INTEGER;
COMMIT ;