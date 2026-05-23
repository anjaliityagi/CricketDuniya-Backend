BEGIN;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'match_teams'
        ) AND EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'match_teams'
              AND column_name = 'original_team_id'
        ) THEN
            ALTER TABLE match_teams
                RENAME COLUMN original_team_id TO team_id;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'match_teams'
        ) THEN
            ALTER TABLE match_teams
                DROP COLUMN IF EXISTS name,
                DROP COLUMN IF EXISTS is_temporary;
        END IF;
    END $$;

DROP INDEX IF EXISTS idx_match_teams_original_team_id;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'match_teams'
        ) THEN
            CREATE INDEX IF NOT EXISTS idx_match_teams_team_id
                ON match_teams(team_id);
        END IF;
    END $$;

ALTER TABLE teams
    DROP COLUMN IF EXISTS match_id,
    DROP COLUMN IF EXISTS created_by_user_id,
    DROP COLUMN IF EXISTS logo_url;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'teams'
              AND column_name = 'is_match_team'
        ) THEN
            ALTER TABLE teams
                RENAME COLUMN is_match_team TO is_temporary;
        END IF;
    END $$;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'team_players'
              AND column_name = 'user_id'
        ) THEN
            ALTER TABLE team_players
                RENAME COLUMN user_id TO player_id;
        END IF;
    END $$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'team_players_player_id_fkey'
        ) THEN
            ALTER TABLE team_players
                ADD CONSTRAINT team_players_player_id_fkey
                    FOREIGN KEY (player_id) REFERENCES users(id);
        END IF;
    END $$;

ALTER TABLE team_players
    DROP COLUMN IF EXISTS is_playing_xi;

DROP INDEX IF EXISTS idx_team_players_user;

CREATE INDEX IF NOT EXISTS idx_team_players_player
    ON team_players(player_id);

COMMIT;
