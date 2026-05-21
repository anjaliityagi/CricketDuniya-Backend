BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";


CREATE TABLE IF NOT EXISTS match_teams
(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id         UUID NOT NULL REFERENCES matches (id) ON DELETE CASCADE,
    source_team_id   UUID REFERENCES teams (id) ON DELETE SET NULL,
    display_name     VARCHAR(100) NOT NULL,
    is_snapshot_edit BOOLEAN          DEFAULT FALSE,
    created_at       TIMESTAMP        DEFAULT NOW(),
    updated_at       TIMESTAMP        DEFAULT NOW(),
    deleted_at       TIMESTAMP,
    UNIQUE (match_id, display_name)
);


CREATE TABLE IF NOT EXISTS match_team_players
(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_team_id    UUID NOT NULL REFERENCES match_teams (id) ON DELETE CASCADE,
    user_id          UUID REFERENCES users (id) ON DELETE SET NULL,
    player_name      VARCHAR(100) NOT NULL,
    phone_number     VARCHAR(20),
    is_playing_xi    BOOLEAN          DEFAULT TRUE,
    is_substitute    BOOLEAN          DEFAULT FALSE,
    is_captain       BOOLEAN          DEFAULT FALSE,
    is_wicket_keeper BOOLEAN          DEFAULT FALSE,
    batting_order    INT,
    created_at       TIMESTAMP        DEFAULT NOW(),
    updated_at       TIMESTAMP        DEFAULT NOW(),
    removed_at       TIMESTAMP,
    deleted_at       TIMESTAMP,
    UNIQUE (match_team_id, user_id)
);


DROP INDEX IF EXISTS uq_match_team_players_captain;
CREATE UNIQUE INDEX uq_match_team_players_captain
    ON match_team_players (match_team_id)
    WHERE is_captain = TRUE AND deleted_at IS NULL;


ALTER TABLE matches
    ADD COLUMN IF NOT EXISTS scorer_user_id UUID REFERENCES users (id),
    ADD COLUMN IF NOT EXISTS toss_winner_match_team_id UUID REFERENCES match_teams (id),
    ADD COLUMN IF NOT EXISTS winner_match_team_id UUID REFERENCES match_teams (id),
    ADD COLUMN IF NOT EXISTS player_of_match_user_id UUID REFERENCES users (id),
    ADD COLUMN IF NOT EXISTS worst_player_user_id UUID REFERENCES users (id),
    ADD COLUMN IF NOT EXISTS started_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP;

ALTER TABLE innings
    ADD COLUMN IF NOT EXISTS batting_match_team_id UUID REFERENCES match_teams (id),
    ADD COLUMN IF NOT EXISTS bowling_match_team_id UUID REFERENCES match_teams (id);

ALTER TABLE ball_events
    ADD COLUMN IF NOT EXISTS batting_match_team_id UUID REFERENCES match_teams (id),
    ADD COLUMN IF NOT EXISTS bowling_match_team_id UUID REFERENCES match_teams (id),
    ADD COLUMN IF NOT EXISTS striker_match_player_id UUID REFERENCES match_team_players (id),
    ADD COLUMN IF NOT EXISTS non_striker_match_player_id UUID REFERENCES match_team_players (id),
    ADD COLUMN IF NOT EXISTS bowler_match_player_id UUID REFERENCES match_team_players (id),
    ADD COLUMN IF NOT EXISTS dismissed_match_player_id UUID REFERENCES match_team_players (id),
    ADD COLUMN IF NOT EXISTS dismissal_type VARCHAR(50),
    ADD COLUMN IF NOT EXISTS is_legal_delivery BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS scoring_version VARCHAR(20) DEFAULT 'v1';

UPDATE ball_events
SET striker_match_player_id = striker_id
WHERE striker_match_player_id IS NULL
  AND striker_id IS NOT NULL;

UPDATE ball_events
SET non_striker_match_player_id = non_striker_id
WHERE non_striker_match_player_id IS NULL
  AND non_striker_id IS NOT NULL;

UPDATE ball_events
SET bowler_match_player_id = bowler_id
WHERE bowler_match_player_id IS NULL
  AND bowler_id IS NOT NULL;

UPDATE ball_events
SET dismissed_match_player_id = dismissed_player_id
WHERE dismissed_match_player_id IS NULL
  AND dismissed_player_id IS NOT NULL;


ALTER TABLE player_match_stats
    ADD COLUMN IF NOT EXISTS match_team_player_id UUID REFERENCES match_team_players (id),
    ADD COLUMN IF NOT EXISTS batting_points INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS bowling_points INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS fielding_points INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS result_points INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();




DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'point_category') THEN
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
    END
$$;

CREATE TABLE IF NOT EXISTS point_events
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id      UUID NOT NULL REFERENCES matches (id) ON DELETE CASCADE,
    user_id       UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    ball_event_id UUID REFERENCES ball_events (id) ON DELETE CASCADE,
    category      point_category NOT NULL,
    rule_name     VARCHAR(120)   NOT NULL,
    points        INT            NOT NULL,
    created_at    TIMESTAMP               DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_match_teams_match_id ON match_teams (match_id);
CREATE INDEX IF NOT EXISTS idx_match_team_players_match_team_id ON match_team_players (match_team_id);
CREATE INDEX IF NOT EXISTS idx_match_team_players_user_id ON match_team_players (user_id);
CREATE INDEX IF NOT EXISTS idx_ball_events_match_player_refs
    ON ball_events (match_id, innings_id, over_no, ball_no, striker_match_player_id, bowler_match_player_id);
CREATE INDEX IF NOT EXISTS idx_player_match_stats_match_player
    ON player_match_stats (match_id, match_team_player_id);
CREATE INDEX IF NOT EXISTS idx_point_events_match_user
    ON point_events (match_id, user_id);

COMMIT;
