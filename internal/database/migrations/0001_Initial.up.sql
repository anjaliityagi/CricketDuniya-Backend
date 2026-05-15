
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE toss_decision AS ENUM ('bat', 'bowl');

CREATE TYPE tournament_status AS ENUM ('upcoming', 'ongoing', 'completed');

CREATE TYPE match_status AS ENUM ('scheduled', 'live', 'completed', 'cancelled');

CREATE TYPE extra_type AS ENUM ('wide', 'no_ball');

CREATE TYPE wicket_type AS ENUM (
    'bowled',
    'caught',
    'lbw',
    'run_out',
    'stumped',
    'hit_wicket'
);


CREATE TABLE users (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(100) NOT NULL,
email VARCHAR(100) UNIQUE NOT NULL,
 password_hash TEXT NOT NULL,
profile_image TEXT,
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE teams (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 name VARCHAR(100) NOT NULL,
 logo_url TEXT,
captain_id UUID REFERENCES users(id),
 created_by UUID REFERENCES users(id),
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE players (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     user_id UUID REFERENCES users(id),
     team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
     full_name VARCHAR(100) NOT NULL,
     is_captain BOOLEAN DEFAULT FALSE,
    is_guest BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE tournaments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     name VARCHAR(150) NOT NULL,
     location TEXT,
     start_date DATE,
     end_date DATE,
  organizer_id UUID REFERENCES users(id),
status tournament_status DEFAULT 'upcoming',
created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tournament_teams (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
tournament_id UUID REFERENCES tournaments(id) ON DELETE CASCADE,
team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
points INT DEFAULT 0,
 matches_played INT DEFAULT 0,
 wins INT DEFAULT 0,
 losses INT DEFAULT 0,
created_at TIMESTAMP DEFAULT NOW(),
UNIQUE(tournament_id, team_id)
);


CREATE TABLE matches (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tournament_id UUID REFERENCES tournaments(id),
    team_a_id UUID REFERENCES teams(id),
    team_b_id UUID REFERENCES teams(id),
     host_user_id UUID NOT NULL REFERENCES users(id),
    venue TEXT,
    match_date TIMESTAMP,
    overs_per_side INT NOT NULL,
     toss_winner_team_id UUID REFERENCES teams(id),
    toss_decision toss_decision,
    winner_team_id UUID REFERENCES teams(id),
    status match_status DEFAULT 'scheduled',
 created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE innings (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
  batting_team_id UUID REFERENCES teams(id),
    bowling_team_id UUID REFERENCES teams(id),
  innings_number INT NOT NULL,
     total_runs INT DEFAULT 0,
    wickets INT DEFAULT 0,
    total_overs DECIMAL(4,1) DEFAULT 0,
    started_at TIMESTAMP,
  ended_at TIMESTAMP
);


CREATE TABLE overs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    innings_id UUID REFERENCES innings(id) ON DELETE CASCADE,
     over_number INT NOT NULL,
    bowler_id UUID REFERENCES players(id),
    runs_conceded INT DEFAULT 0,
      wickets INT DEFAULT 0,
  completed BOOLEAN DEFAULT FALSE
);
CREATE TABLE deliveries
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    innings_id     UUID REFERENCES innings (id) ON DELETE CASCADE,
    over_id        UUID REFERENCES overs (id) ON DELETE CASCADE,
    striker_id     UUID REFERENCES players (id),
    non_striker_id UUID REFERENCES players (id),
    bowler_id      UUID REFERENCES players (id),
    ball_number    INT NOT NULL,
    runs_scored    INT              DEFAULT 0,
    extras         INT              DEFAULT 0,
    extra_type     extra_type,
    wicket         BOOLEAN          DEFAULT FALSE,
    wicket_type    wicket_type,
    player_out_id  UUID REFERENCES players (id),
    created_at     TIMESTAMP        DEFAULT NOW()
);
CREATE TABLE super_overs
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id        UUID NOT NULL REFERENCES matches (id) ON DELETE CASCADE,
    team1_id        UUID NOT NULL REFERENCES teams (id),
    team2_id        UUID NOT NULL REFERENCES teams (id),
    batting_team_id UUID NOT NULL REFERENCES teams (id),
    bowling_team_id UUID NOT NULL REFERENCES teams (id),
    created_at      TIMESTAMPTZ      DEFAULT NOW(),
    updated_at      TIMESTAMPTZ      DEFAULT NOW()
);