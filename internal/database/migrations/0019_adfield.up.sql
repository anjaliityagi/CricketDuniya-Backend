BEGIN ;
ALTER TABLE team_players
    ADD COLUMN is_playing_xi BOOLEAN DEFAULT true;

COMMIT ;