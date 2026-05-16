BEGIN;
DROP TABLE IF EXISTS players CASCADE ;
CREATE TABLE IF NOT EXISTS match_players
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id        UUID         NOT NULL REFERENCES matches (id) ON DELETE CASCADE,
    team_id         UUID         NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    user_id         UUID         REFERENCES users (id) ON DELETE SET NULL,
    player_name     VARCHAR(100) NOT NULL,
    phone           VARCHAR(20),
    is_captain      BOOLEAN          DEFAULT FALSE,
    is_wicketkeeper BOOLEAN          DEFAULT FALSE,
    created_at      TIMESTAMP        DEFAULT NOW()
);
DROP TABLE match_players CASCADE ;

CREATE TABLE IF NOT EXISTS match_players
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,

    user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    player_name VARCHAR(100) NOT NULL,

    phone VARCHAR(20),

    team_side VARCHAR(1),

    is_host BOOLEAN DEFAULT FALSE,

    is_captain BOOLEAN DEFAULT FALSE,

    is_wicketkeeper BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT NOW()
);
COMMIT;