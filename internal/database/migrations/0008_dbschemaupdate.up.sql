
BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =====================================================
-- DROP UNUSED TABLES
-- =====================================================

DROP TABLE IF EXISTS tournament_teams CASCADE;
DROP TABLE IF EXISTS tournaments CASCADE;
DROP TABLE IF EXISTS overs CASCADE;
DROP TABLE IF EXISTS players CASCADE;

-- =====================================================
-- USERS
-- =====================================================

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_phone_verified BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

ALTER TABLE users
    DROP COLUMN IF EXISTS batting_style,
    DROP COLUMN IF EXISTS bowling_style,
    DROP COLUMN IF EXISTS role;

-- =====================================================
-- TEAMS
-- =====================================================

ALTER TABLE teams
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

ALTER TABLE teams
    RENAME COLUMN created_by TO created_by_user_id;

-- =====================================================
-- TEAM_PLAYERS
-- =====================================================

ALTER TABLE team_players
    RENAME COLUMN player_id TO user_id;

ALTER TABLE team_players
    ADD COLUMN IF NOT EXISTS batting_style VARCHAR(50),
    ADD COLUMN IF NOT EXISTS bowling_style VARCHAR(50),
    ADD COLUMN IF NOT EXISTS role VARCHAR(50),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

ALTER TABLE team_players
    ADD CONSTRAINT unique_team_user UNIQUE(team_id, user_id);

-- =====================================================
-- MATCHES
-- =====================================================

ALTER TABLE matches
    DROP COLUMN IF EXISTS tournament_id,
    DROP COLUMN IF EXISTS team_a_id,
    DROP COLUMN IF EXISTS team_b_id,
    DROP COLUMN IF EXISTS current_innings,
    DROP COLUMN IF EXISTS toss_completed,
    DROP COLUMN IF EXISTS winning_margin_runs,
    DROP COLUMN IF EXISTS winning_margin_wickets,
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS total_overs,
    DROP COLUMN IF EXISTS winning_team_id,
    DROP COLUMN IF EXISTS match_status,
    DROP COLUMN IF EXISTS created_by;

ALTER TABLE matches
    RENAME COLUMN venue TO location;

ALTER TABLE matches
    RENAME COLUMN overs_per_side TO overs_per_innings;

ALTER TABLE matches
    RENAME COLUMN match_date TO started_at;

ALTER TABLE matches
    RENAME COLUMN match_ended_at TO completed_at;

ALTER TABLE matches
    RENAME COLUMN winner_team_id TO winner_match_team_id;

ALTER TABLE matches
    ADD COLUMN IF NOT EXISTS scorer_user_id UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS current_innings_no INT DEFAULT 1,
    ADD COLUMN IF NOT EXISTS player_of_match_user_id UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS worst_player_user_id UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- =====================================================
-- MATCH_TEAMS
-- =====================================================

ALTER TABLE match_teams
    DROP COLUMN IF EXISTS target_score;

ALTER TABLE match_teams
    RENAME COLUMN source_team_id TO original_team_id;

ALTER TABLE match_teams
    ADD COLUMN IF NOT EXISTS is_temporary BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- =====================================================
-- REMOVE MATCH_TEAM_PLAYERS
-- =====================================================

DROP TABLE IF EXISTS match_team_players CASCADE;

-- =====================================================
-- INNINGS
-- =====================================================

ALTER TABLE innings
    DROP COLUMN IF EXISTS batting_team_id,
    DROP COLUMN IF EXISTS bowling_team_id,
    DROP COLUMN IF EXISTS total_overs,
    DROP COLUMN IF EXISTS target_runs,
    DROP COLUMN IF EXISTS is_completed;

ALTER TABLE innings
    RENAME COLUMN innings_number TO innings_no;

ALTER TABLE innings
    RENAME COLUMN wickets TO total_wickets;

ALTER TABLE innings
    RENAME COLUMN ended_at TO completed_at;

-- =====================================================
-- BALLS -> BALL_EVENTS
-- =====================================================

ALTER TABLE balls
    RENAME TO ball_events;

ALTER TABLE ball_events
    DROP COLUMN IF EXISTS over_id,
    DROP COLUMN IF EXISTS extra_type,
    DROP COLUMN IF EXISTS wicket,
    DROP COLUMN IF EXISTS wicket_type,
    DROP COLUMN IF EXISTS is_boundary,
    DROP COLUMN IF EXISTS is_six,
    DROP COLUMN IF EXISTS total_runs,
    DROP COLUMN IF EXISTS is_valid_ball,
    DROP COLUMN IF EXISTS is_void,
    DROP COLUMN IF EXISTS void_reason;

ALTER TABLE ball_events
    RENAME COLUMN striker_id TO striker_team_player_id;

ALTER TABLE ball_events
    RENAME COLUMN non_striker_id TO non_striker_team_player_id;

ALTER TABLE ball_events
    RENAME COLUMN bowler_id TO bowler_team_player_id;

ALTER TABLE ball_events
    RENAME COLUMN player_out_id TO dismissed_team_player_id;

ALTER TABLE ball_events
    RENAME COLUMN fielder_id TO fielder_user_id;

ALTER TABLE ball_events
    RENAME COLUMN ball_number TO ball_no;

ALTER TABLE ball_events
    ADD COLUMN IF NOT EXISTS match_id UUID REFERENCES matches(id),
    ADD COLUMN IF NOT EXISTS over_no INT DEFAULT 1,
    ADD COLUMN IF NOT EXISTS delivery_no INT DEFAULT 1,
    ADD COLUMN IF NOT EXISTS is_legal_delivery BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS batting_match_team_id UUID REFERENCES match_teams(id),
    ADD COLUMN IF NOT EXISTS bowling_match_team_id UUID REFERENCES match_teams(id),
    ADD COLUMN IF NOT EXISTS wides INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS no_balls INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS byes INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS leg_byes INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS dismissal_type VARCHAR(50),
    ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_by_user_id UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS created_by_user_id UUID REFERENCES users(id);

-- =====================================================
-- PLAYER MATCH STATS
-- =====================================================

ALTER TABLE player_match_stats
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS dismissal_type,
    DROP COLUMN IF EXISTS player_id;

ALTER TABLE player_match_stats
    RENAME COLUMN runs TO runs_scored;

ALTER TABLE player_match_stats
    RENAME COLUMN wickets TO wickets_taken;

ALTER TABLE player_match_stats
    RENAME COLUMN run_outs TO runouts;

ALTER TABLE player_match_stats
    ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS team_player_id UUID REFERENCES team_players(id),
    ADD COLUMN IF NOT EXISTS batting_style VARCHAR(50),
    ADD COLUMN IF NOT EXISTS bowling_style VARCHAR(50),
    ADD COLUMN IF NOT EXISTS batting_points INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS bowling_points INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS fielding_points INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- =====================================================
-- SUPER OVERS
-- =====================================================

CREATE TABLE IF NOT EXISTS super_overs (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                           match_id UUID NOT NULL REFERENCES matches(id),

                                           innings_id UUID REFERENCES innings(id),

                                           batting_match_team_id UUID REFERENCES match_teams(id),

                                           bowling_match_team_id UUID REFERENCES match_teams(id),

                                           winner_match_team_id UUID REFERENCES match_teams(id),

                                           super_over_no INT DEFAULT 1,

                                           created_at TIMESTAMP DEFAULT NOW(),

                                           updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- POINT CATEGORY ENUM
-- =====================================================

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM pg_type WHERE typname = 'point_category'
        ) THEN
            CREATE TYPE point_category AS ENUM (
                'batting',
                'bowling',
                'fielding',
                'milestone',
                'economy',
                'strike_rate',
                'result',
                'penalty'
                );
        END IF;
    END$$;

-- =====================================================
-- POINT EVENTS
-- =====================================================

CREATE TABLE IF NOT EXISTS point_events (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                            match_id UUID NOT NULL REFERENCES matches(id),

                                            user_id UUID NOT NULL REFERENCES users(id),

                                            ball_event_id UUID REFERENCES ball_events(id),

                                            category point_category NOT NULL,

                                            rule_name VARCHAR(100) NOT NULL,

                                            points INT NOT NULL,

                                            created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- AUDIT LOGS
-- =====================================================

CREATE TABLE IF NOT EXISTS audit_logs (
                                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                          user_id UUID REFERENCES users(id),

                                          entity_type VARCHAR(50),

                                          entity_id UUID,

                                          action VARCHAR(50),

                                          old_data JSONB,

                                          new_data JSONB,

                                          created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- INDEXES
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_ball_events_match
    ON ball_events(match_id);

CREATE INDEX IF NOT EXISTS idx_ball_events_over
    ON ball_events(innings_id, over_no, ball_no);

CREATE INDEX IF NOT EXISTS idx_player_match_stats_user
    ON player_match_stats(user_id);

CREATE INDEX IF NOT EXISTS idx_point_events_user
    ON point_events(user_id);

CREATE INDEX IF NOT EXISTS idx_matches_status
    ON matches(status);

COMMIT;
```

# Final Architecture

## Tables Remaining

* users
* user_session
* teams
* team_players
* matches
* match_teams
* innings
* ball_events
* player_match_stats
* point_events
* audit_logs
* super_overs
* schema_migrations

## Tables Removed

* players
* overs
* tournaments
* tournament_teams
