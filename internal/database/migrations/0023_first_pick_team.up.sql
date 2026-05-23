BEGIN;

ALTER TABLE matches
    ADD COLUMN IF NOT EXISTS first_pick_team_id UUID REFERENCES teams (id);

COMMIT;
