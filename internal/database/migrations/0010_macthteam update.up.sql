BEGIN;

ALTER TABLE teams
    ADD COLUMN IF NOT EXISTS match_id UUID REFERENCES matches(id),
    ADD COLUMN IF NOT EXISTS original_team_id UUID REFERENCES teams(id),
    ADD COLUMN IF NOT EXISTS is_match_team BOOLEAN DEFAULT FALSE;

ALTER TABLE team_players
    ADD COLUMN IF NOT EXISTS batting_order INT,
    ADD COLUMN IF NOT EXISTS is_wicket_keeper BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_playing_xi BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS is_substitute BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_guest BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS removed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

ALTER TABLE team_players
    DROP CONSTRAINT IF EXISTS unique_team_user;

ALTER TABLE team_players
    ADD CONSTRAINT unique_team_user_match
        UNIQUE(team_id, user_id);

INSERT INTO team_players (
    id,
    team_id,
    user_id,
    is_captain,
    batting_order,
    is_wicket_keeper,
    is_playing_xi,
    is_substitute,
    is_guest,
    removed_at,
    created_at,
    updated_at,
    deleted_at
)
SELECT
    mtp.id,
    mtp.match_team_id,
    mtp.user_id,
    mtp.is_captain,
    mtp.batting_order,
    mtp.is_wicket_keeper,
    mtp.is_playing_xi,
    mtp.is_substitute,
    mtp.is_guest,
    mtp.removed_at,
    mtp.created_at,
    mtp.updated_at,
    mtp.deleted_at
FROM match_team_players mtp
ON CONFLICT DO NOTHING;


ALTER TABLE ball_events
    RENAME COLUMN striker_match_player_id TO striker_team_player_id;

ALTER TABLE ball_events
    RENAME COLUMN non_striker_match_player_id TO non_striker_team_player_id;

ALTER TABLE ball_events
    RENAME COLUMN bowler_match_player_id TO bowler_team_player_id;

ALTER TABLE ball_events
    RENAME COLUMN dismissed_match_player_id TO dismissed_team_player_id;



ALTER TABLE player_match_stats
    RENAME COLUMN match_team_player_id TO team_player_id;


DROP TABLE IF EXISTS match_team_players CASCADE;

CREATE INDEX IF NOT EXISTS idx_team_players_team
    ON team_players(team_id);

CREATE INDEX IF NOT EXISTS idx_team_players_user
    ON team_players(user_id);

CREATE INDEX IF NOT EXISTS idx_teams_match
    ON teams(match_id);

COMMIT;