BEGIN;

ALTER TABLE match_players
    ADD COLUMN IF NOT EXISTS removed_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_match_players_removed_at
    ON match_players(removed_at);

COMMIT;
