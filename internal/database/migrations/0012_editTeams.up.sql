BEGIN ;
ALTER TABLE  if exists match_teams
    RENAME COLUMN original_team_id TO team_id;

ALTER TABLE if exists match_teams
    DROP COLUMN name;

ALTER TABLE if exists match_teams
    DROP COLUMN is_temporary;

DROP INDEX IF EXISTS idx_match_teams_original_team_id;

CREATE INDEX IF NOT EXISTS idx_match_teams_team_id
    ON match_teams(team_id);

ALTER TABLE if exists teams
    DROP COLUMN match_id;

ALTER TABLE if exists teams
    DROP COLUMN created_by_user_id;

ALTER TABLE if exists teams
    DROP COLUMN logo_url;

ALTER TABLE if exists teams
    RENAME COLUMN is_match_team TO is_temporary;


ALTER TABLE if exists team_players
    RENAME COLUMN user_id TO player_id;

ALTER TABLE if exists team_players
    ADD CONSTRAINT team_players_player_id_fkey
        FOREIGN KEY (player_id) REFERENCES users(id);

ALTER TABLE if exists team_players
    DROP COLUMN is_playing_xi;

DROP INDEX IF EXISTS idx_team_players_user;
CREATE INDEX IF NOT EXISTS idx_team_players_player
    ON team_players(player_id);

COMMIT ;