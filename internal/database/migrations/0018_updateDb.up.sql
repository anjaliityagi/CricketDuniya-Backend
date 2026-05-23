BEGIN;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'match_teams'
        ) THEN
            IF EXISTS (
                SELECT 1 FROM information_schema.columns
                WHERE table_name = 'innings'
                  AND column_name = 'batting_match_team_id'
            ) THEN
                UPDATE innings i
                SET batting_match_team_id = mt.source_team_id
                FROM match_teams mt
                WHERE i.batting_match_team_id = mt.id
                  AND mt.source_team_id IS NOT NULL;
            END IF;

            IF EXISTS (
                SELECT 1 FROM information_schema.columns
                WHERE table_name = 'innings'
                  AND column_name = 'bowling_match_team_id'
            ) THEN
                UPDATE innings i
                SET bowling_match_team_id = mt.source_team_id
                FROM match_teams mt
                WHERE i.bowling_match_team_id = mt.id
                  AND mt.source_team_id IS NOT NULL;
            END IF;
        END IF;
    END $$;

DROP TABLE IF EXISTS match_team_players CASCADE;
DROP TABLE IF EXISTS match_teams CASCADE;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'ball_events'
        ) THEN
            ALTER TABLE ball_events
                DROP COLUMN IF EXISTS striker_match_player_id,
                DROP COLUMN IF EXISTS non_striker_match_player_id,
                DROP COLUMN IF EXISTS bowler_match_player_id,
                DROP COLUMN IF EXISTS dismissed_match_player_id,
                DROP COLUMN IF EXISTS batting_match_team_id,
                DROP COLUMN IF EXISTS bowling_match_team_id;
        END IF;
    END $$;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'innings'
              AND column_name = 'batting_match_team_id'
        ) THEN
            ALTER TABLE innings
                RENAME COLUMN batting_match_team_id TO batting_team_id;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'innings'
              AND column_name = 'bowling_match_team_id'
        ) THEN
            ALTER TABLE innings
                RENAME COLUMN bowling_match_team_id TO bowling_team_id;
        END IF;
    END $$;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'player_match_stats'
              AND column_name = 'match_team_player_id'
        ) AND NOT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'player_match_stats'
              AND column_name = 'team_player_id'
        ) THEN
            ALTER TABLE player_match_stats
                RENAME COLUMN match_team_player_id TO team_player_id;
        END IF;
    END $$;

ALTER TABLE player_match_stats
    DROP COLUMN IF EXISTS match_team_player_id;

ALTER TABLE player_match_stats
    ADD COLUMN IF NOT EXISTS team_player_id UUID;

ALTER TABLE teams
    DROP COLUMN IF EXISTS original_team_id,
    DROP COLUMN IF EXISTS is_temporary;

ALTER TABLE team_players
    DROP COLUMN IF EXISTS removed_at;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'matches'
              AND column_name = 'team_a_id'
        ) THEN
            UPDATE matches
            SET team_a_id = NULL
            WHERE team_a_id IS NOT NULL
              AND NOT EXISTS (SELECT 1 FROM teams t WHERE t.id = matches.team_a_id);
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'matches'
              AND column_name = 'team_b_id'
        ) THEN
            UPDATE matches
            SET team_b_id = NULL
            WHERE team_b_id IS NOT NULL
              AND NOT EXISTS (SELECT 1 FROM teams t WHERE t.id = matches.team_b_id);
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'innings'
              AND column_name = 'batting_team_id'
        ) THEN
            UPDATE innings
            SET batting_team_id = NULL
            WHERE batting_team_id IS NOT NULL
              AND NOT EXISTS (SELECT 1 FROM teams t WHERE t.id = innings.batting_team_id);
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'innings'
              AND column_name = 'bowling_team_id'
        ) THEN
            UPDATE innings
            SET bowling_team_id = NULL
            WHERE bowling_team_id IS NOT NULL
              AND NOT EXISTS (SELECT 1 FROM teams t WHERE t.id = innings.bowling_team_id);
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'striker_id'
        ) THEN
            UPDATE ball_events
            SET striker_id = NULL
            WHERE striker_id IS NOT NULL
              AND NOT EXISTS (
                  SELECT 1 FROM team_players tp WHERE tp.id = ball_events.striker_id
              );
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'non_striker_id'
        ) THEN
            UPDATE ball_events
            SET non_striker_id = NULL
            WHERE non_striker_id IS NOT NULL
              AND NOT EXISTS (
                  SELECT 1 FROM team_players tp WHERE tp.id = ball_events.non_striker_id
              );
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'bowler_id'
        ) THEN
            UPDATE ball_events
            SET bowler_id = NULL
            WHERE bowler_id IS NOT NULL
              AND NOT EXISTS (
                  SELECT 1 FROM team_players tp WHERE tp.id = ball_events.bowler_id
              );
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'dismissed_player_id'
        ) THEN
            UPDATE ball_events
            SET dismissed_player_id = NULL
            WHERE dismissed_player_id IS NOT NULL
              AND NOT EXISTS (
                  SELECT 1 FROM team_players tp WHERE tp.id = ball_events.dismissed_player_id
              );
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'ball_events'
              AND column_name = 'fielder_id'
        ) THEN
            UPDATE ball_events
            SET fielder_id = NULL
            WHERE fielder_id IS NOT NULL
              AND NOT EXISTS (
                  SELECT 1 FROM team_players tp WHERE tp.id = ball_events.fielder_id
              );
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'player_match_stats'
              AND column_name = 'team_player_id'
        ) THEN
            UPDATE player_match_stats
            SET team_player_id = NULL
            WHERE team_player_id IS NOT NULL
              AND NOT EXISTS (
                  SELECT 1 FROM team_players tp WHERE tp.id = player_match_stats.team_player_id
              );
        END IF;
    END $$;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'ball_events'
        ) THEN
            ALTER TABLE ball_events
                DROP CONSTRAINT IF EXISTS deliveries_striker_id_fkey,
                DROP CONSTRAINT IF EXISTS deliveries_non_striker_id_fkey,
                DROP CONSTRAINT IF EXISTS deliveries_bowler_id_fkey,
                DROP CONSTRAINT IF EXISTS deliveries_player_out_id_fkey;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'matches'
        ) THEN
            ALTER TABLE matches
                DROP CONSTRAINT IF EXISTS matches_team_a_id_fkey,
                DROP CONSTRAINT IF EXISTS matches_team_b_id_fkey,
                DROP CONSTRAINT IF EXISTS fk_matches_team_a,
                DROP CONSTRAINT IF EXISTS fk_matches_team_b;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'innings'
        ) THEN
            ALTER TABLE innings
                DROP CONSTRAINT IF EXISTS innings_batting_team_id_fkey,
                DROP CONSTRAINT IF EXISTS innings_bowling_team_id_fkey,
                DROP CONSTRAINT IF EXISTS fk_innings_batting_team,
                DROP CONSTRAINT IF EXISTS fk_innings_bowling_team;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'ball_events'
        ) THEN
            ALTER TABLE ball_events
                DROP CONSTRAINT IF EXISTS fk_ball_events_striker,
                DROP CONSTRAINT IF EXISTS fk_ball_events_non_striker,
                DROP CONSTRAINT IF EXISTS fk_ball_events_bowler,
                DROP CONSTRAINT IF EXISTS fk_ball_events_dismissed,
                DROP CONSTRAINT IF EXISTS fk_ball_events_fielder;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'player_match_stats'
        ) THEN
            ALTER TABLE player_match_stats
                DROP CONSTRAINT IF EXISTS fk_player_match_stats_team_player,
                DROP CONSTRAINT IF EXISTS player_match_stats_player_id_fkey;
        END IF;
    END $$;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'matches'
              AND column_name = 'team_a_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'matches'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'team_a_id'
        ) THEN
            ALTER TABLE matches
                ADD CONSTRAINT fk_matches_team_a
                    FOREIGN KEY (team_a_id)
                        REFERENCES teams(id)
                        ON DELETE CASCADE;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'matches'
              AND column_name = 'team_b_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'matches'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'team_b_id'
        ) THEN
            ALTER TABLE matches
                ADD CONSTRAINT fk_matches_team_b
                    FOREIGN KEY (team_b_id)
                        REFERENCES teams(id)
                        ON DELETE CASCADE;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'innings'
              AND column_name = 'batting_team_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'innings'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'batting_team_id'
        ) THEN
            ALTER TABLE innings
                ADD CONSTRAINT fk_innings_batting_team
                    FOREIGN KEY (batting_team_id)
                        REFERENCES teams(id)
                        ON DELETE CASCADE;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'innings'
              AND column_name = 'bowling_team_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'innings'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'bowling_team_id'
        ) THEN
            ALTER TABLE innings
                ADD CONSTRAINT fk_innings_bowling_team
                    FOREIGN KEY (bowling_team_id)
                        REFERENCES teams(id)
                        ON DELETE CASCADE;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'ball_events'
              AND column_name = 'striker_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'ball_events'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'striker_id'
        ) THEN
            ALTER TABLE ball_events
                ADD CONSTRAINT fk_ball_events_striker
                    FOREIGN KEY (striker_id)
                        REFERENCES team_players(id)
                        ON DELETE SET NULL;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'ball_events'
              AND column_name = 'non_striker_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'ball_events'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'non_striker_id'
        ) THEN
            ALTER TABLE ball_events
                ADD CONSTRAINT fk_ball_events_non_striker
                    FOREIGN KEY (non_striker_id)
                        REFERENCES team_players(id)
                        ON DELETE SET NULL;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'ball_events'
              AND column_name = 'bowler_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'ball_events'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'bowler_id'
        ) THEN
            ALTER TABLE ball_events
                ADD CONSTRAINT fk_ball_events_bowler
                    FOREIGN KEY (bowler_id)
                        REFERENCES team_players(id)
                        ON DELETE SET NULL;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'ball_events'
              AND column_name = 'dismissed_player_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'ball_events'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'dismissed_player_id'
        ) THEN
            ALTER TABLE ball_events
                ADD CONSTRAINT fk_ball_events_dismissed
                    FOREIGN KEY (dismissed_player_id)
                        REFERENCES team_players(id)
                        ON DELETE SET NULL;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'ball_events'
              AND column_name = 'fielder_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'ball_events'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'fielder_id'
        ) THEN
            ALTER TABLE ball_events
                ADD CONSTRAINT fk_ball_events_fielder
                    FOREIGN KEY (fielder_id)
                        REFERENCES team_players(id)
                        ON DELETE SET NULL;
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'player_match_stats'
              AND column_name = 'team_player_id'
        ) AND NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_schema = kcu.constraint_schema
               AND tc.constraint_name = kcu.constraint_name
            WHERE tc.table_schema = 'public'
              AND tc.table_name = 'player_match_stats'
              AND tc.constraint_type = 'FOREIGN KEY'
              AND kcu.column_name = 'team_player_id'
        ) THEN
            ALTER TABLE player_match_stats
                ADD CONSTRAINT fk_player_match_stats_team_player
                    FOREIGN KEY (team_player_id)
                        REFERENCES team_players(id)
                        ON DELETE CASCADE;
        END IF;
    END $$;

CREATE INDEX IF NOT EXISTS idx_ball_events_match_id
    ON ball_events(match_id);

CREATE INDEX IF NOT EXISTS idx_ball_events_innings_id
    ON ball_events(innings_id);

CREATE INDEX IF NOT EXISTS idx_team_players_team_id
    ON team_players(team_id);

CREATE INDEX IF NOT EXISTS idx_player_match_stats_match_id
    ON player_match_stats(match_id);

CREATE INDEX IF NOT EXISTS idx_player_match_stats_team_player_id
    ON player_match_stats(team_player_id);

COMMIT;
