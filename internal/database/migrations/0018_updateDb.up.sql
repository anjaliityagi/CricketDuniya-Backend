BEGIN ;



DROP TABLE IF EXISTS match_team_players CASCADE;
DROP TABLE IF EXISTS match_teams CASCADE;



ALTER TABLE ball_events
    DROP COLUMN IF EXISTS striker_match_player_id,
    DROP COLUMN IF EXISTS non_striker_match_player_id,
    DROP COLUMN IF EXISTS bowler_match_player_id,
    DROP COLUMN IF EXISTS dismissed_match_player_id,
    DROP COLUMN IF EXISTS batting_match_team_id,
    DROP COLUMN IF EXISTS bowling_match_team_id;



ALTER TABLE innings
    RENAME COLUMN batting_match_team_id TO batting_team_id;

ALTER TABLE innings
    RENAME COLUMN bowling_match_team_id TO bowling_team_id;



ALTER TABLE player_match_stats
    DROP COLUMN IF EXISTS match_team_player_id;



ALTER TABLE teams
    DROP COLUMN IF EXISTS original_team_id,
    DROP COLUMN IF EXISTS is_temporary;



ALTER TABLE team_players
    DROP COLUMN IF EXISTS removed_at;


ALTER TABLE matches
    ADD CONSTRAINT fk_matches_team_a
        FOREIGN KEY (team_a_id)
            REFERENCES teams(id)
            ON DELETE CASCADE;

ALTER TABLE matches
    ADD CONSTRAINT fk_matches_team_b
        FOREIGN KEY (team_b_id)
            REFERENCES teams(id)
            ON DELETE CASCADE;


ALTER TABLE innings
    ADD CONSTRAINT fk_innings_batting_team
        FOREIGN KEY (batting_team_id)
            REFERENCES teams(id)
            ON DELETE CASCADE;

ALTER TABLE innings
    ADD CONSTRAINT fk_innings_bowling_team
        FOREIGN KEY (bowling_team_id)
            REFERENCES teams(id)
            ON DELETE CASCADE;

ALTER TABLE ball_events
    ADD CONSTRAINT fk_ball_events_striker
        FOREIGN KEY (striker_id)
            REFERENCES team_players(id)
            ON DELETE SET NULL;

ALTER TABLE ball_events
    ADD CONSTRAINT fk_ball_events_non_striker
        FOREIGN KEY (non_striker_id)
            REFERENCES team_players(id)
            ON DELETE SET NULL;

ALTER TABLE ball_events
    ADD CONSTRAINT fk_ball_events_bowler
        FOREIGN KEY (bowler_id)
            REFERENCES team_players(id)
            ON DELETE SET NULL;

ALTER TABLE ball_events
    ADD CONSTRAINT fk_ball_events_dismissed
        FOREIGN KEY (dismissed_player_id)
            REFERENCES team_players(id)
            ON DELETE SET NULL;

ALTER TABLE ball_events
    ADD CONSTRAINT fk_ball_events_fielder
        FOREIGN KEY (fielder_id)
            REFERENCES team_players(id)
            ON DELETE SET NULL;

ALTER TABLE player_match_stats
    ADD CONSTRAINT fk_player_match_stats_team_player
        FOREIGN KEY (team_player_id)
            REFERENCES team_players(id)
            ON DELETE CASCADE;



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
COMMIT ;