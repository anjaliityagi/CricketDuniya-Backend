BEGIN;


ALTER TYPE match_status ADD VALUE IF NOT EXISTS 'toss';
ALTER TYPE match_status ADD VALUE IF NOT EXISTS 'innings_break';


ALTER TABLE matches
    ADD COLUMN IF NOT EXISTS current_innings INT DEFAULT 1;

ALTER TABLE matches
    ADD COLUMN IF NOT EXISTS toss_completed BOOLEAN DEFAULT FALSE;

ALTER TABLE matches
    ADD COLUMN IF NOT EXISTS winning_margin_runs INT;

ALTER TABLE matches
    ADD COLUMN IF NOT EXISTS winning_margin_wickets INT;

ALTER TABLE matches
    ADD COLUMN IF NOT EXISTS match_ended_at TIMESTAMP;



ALTER TABLE innings
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'live';

ALTER TABLE innings
    ADD COLUMN IF NOT EXISTS target_runs INT;


ALTER TABLE match_players
    DROP COLUMN IF EXISTS team_side;

ALTER TABLE match_players
    DROP COLUMN IF EXISTS bowling_position;

-- add proper team relation if missing
ALTER TABLE match_players
    ADD COLUMN IF NOT EXISTS team_id UUID
        REFERENCES teams(id)
            ON DELETE CASCADE;

ALTER TABLE match_players
    ADD COLUMN IF NOT EXISTS batting_position INT;

ALTER TABLE match_players
    ADD COLUMN IF NOT EXISTS is_playing_xi BOOLEAN DEFAULT TRUE;

ALTER TABLE match_players
    ADD COLUMN IF NOT EXISTS is_substitute BOOLEAN DEFAULT FALSE;



ALTER TABLE overs
    DROP CONSTRAINT IF EXISTS overs_bowler_id_fkey;


ALTER TABLE overs
    ADD CONSTRAINT overs_bowler_id_fkey
        FOREIGN KEY (bowler_id)
            REFERENCES match_players(id)
            ON DELETE SET NULL;

ALTER TABLE overs
    ADD COLUMN IF NOT EXISTS maiden_over BOOLEAN DEFAULT FALSE;

ALTER TABLE overs
    ADD COLUMN IF NOT EXISTS no_boundary_over BOOLEAN DEFAULT FALSE;

ALTER TABLE deliveries
    DROP CONSTRAINT IF EXISTS deliveries_striker_id_fkey;

ALTER TABLE deliveries
    DROP CONSTRAINT IF EXISTS deliveries_non_striker_id_fkey;

ALTER TABLE deliveries
    DROP CONSTRAINT IF EXISTS deliveries_bowler_id_fkey;

ALTER TABLE deliveries
    DROP CONSTRAINT IF EXISTS deliveries_player_out_id_fkey;


ALTER TABLE deliveries
    ADD CONSTRAINT deliveries_striker_id_fkey
        FOREIGN KEY (striker_id)
            REFERENCES match_players(id)
            ON DELETE SET NULL;

ALTER TABLE deliveries
    ADD CONSTRAINT deliveries_non_striker_id_fkey
        FOREIGN KEY (non_striker_id)
            REFERENCES match_players(id)
            ON DELETE SET NULL;

ALTER TABLE deliveries
    ADD CONSTRAINT deliveries_bowler_id_fkey
        FOREIGN KEY (bowler_id)
            REFERENCES match_players(id)
            ON DELETE SET NULL;

ALTER TABLE deliveries
    ADD CONSTRAINT deliveries_player_out_id_fkey
        FOREIGN KEY (player_out_id)
            REFERENCES match_players(id)
            ON DELETE SET NULL;



ALTER TABLE deliveries
    ADD COLUMN IF NOT EXISTS is_boundary BOOLEAN DEFAULT FALSE;

ALTER TABLE deliveries
    ADD COLUMN IF NOT EXISTS is_six BOOLEAN DEFAULT FALSE;

ALTER TABLE deliveries
    ADD COLUMN IF NOT EXISTS is_dot_ball BOOLEAN DEFAULT FALSE;

ALTER TABLE deliveries
    ADD COLUMN IF NOT EXISTS commentary TEXT;


CREATE TABLE IF NOT EXISTS player_match_stats
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_id UUID NOT NULL
        REFERENCES matches(id)
            ON DELETE CASCADE,

    player_id UUID NOT NULL
        REFERENCES match_players(id)
            ON DELETE CASCADE,


    runs INT DEFAULT 0,
    balls_faced INT DEFAULT 0,
    fours INT DEFAULT 0,
    sixes INT DEFAULT 0,
    strike_rate DECIMAL(6,2) DEFAULT 0,
    is_out BOOLEAN DEFAULT FALSE,


    overs_bowled DECIMAL(4,1) DEFAULT 0,
    maidens INT DEFAULT 0,
    runs_conceded INT DEFAULT 0,
    wickets INT DEFAULT 0,
    economy DECIMAL(6,2) DEFAULT 0,


    catches INT DEFAULT 0,
    stumping INT DEFAULT 0,
    run_outs INT DEFAULT 0,


    fantasy_points INT DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(match_id, player_id)
);


CREATE INDEX IF NOT EXISTS idx_deliveries_innings_id
    ON deliveries(innings_id);

CREATE INDEX IF NOT EXISTS idx_deliveries_over_id
    ON deliveries(over_id);

CREATE INDEX IF NOT EXISTS idx_overs_innings_id
    ON overs(innings_id);

CREATE INDEX IF NOT EXISTS idx_match_players_match_id
    ON match_players(match_id);

CREATE INDEX IF NOT EXISTS idx_player_match_stats_match_id
    ON player_match_stats(match_id);

COMMIT;