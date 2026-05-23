BEGIN;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'striker_team_player_id'
        ) THEN
            ALTER TABLE ball_events
                RENAME COLUMN striker_team_player_id TO striker_id;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'non_striker_team_player_id'
        ) THEN
            ALTER TABLE ball_events
                RENAME COLUMN non_striker_team_player_id TO non_striker_id;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'bowler_team_player_id'
        ) THEN
            ALTER TABLE ball_events
                RENAME COLUMN bowler_team_player_id TO bowler_id;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'dismissed_team_player_id'
        ) THEN
            ALTER TABLE ball_events
                RENAME COLUMN dismissed_team_player_id TO dismissed_player_id;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'fielder_user_id'
        ) THEN
            ALTER TABLE ball_events
                RENAME COLUMN fielder_user_id TO fielder_id;
        END IF;
    END $$;

ALTER TABLE ball_events
    ADD COLUMN IF NOT EXISTS delivery_number INTEGER,
    ADD COLUMN IF NOT EXISTS ball_type VARCHAR(20),
    ADD COLUMN IF NOT EXISTS total_runs INTEGER,
    ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_by_user_id UUID,
    ADD COLUMN IF NOT EXISTS created_by_user_id UUID;

CREATE INDEX IF NOT EXISTS idx_ball_events_match
    ON ball_events(match_id);

CREATE INDEX IF NOT EXISTS idx_ball_events_innings
    ON ball_events(innings_id);

CREATE INDEX IF NOT EXISTS idx_ball_events_over_ball
    ON ball_events(over_no, ball_no);

COMMIT;
