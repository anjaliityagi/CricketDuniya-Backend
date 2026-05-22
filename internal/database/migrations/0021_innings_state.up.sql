BEGIN;

CREATE TABLE IF NOT EXISTS innings_state (
    innings_id UUID PRIMARY KEY REFERENCES innings(id) ON DELETE CASCADE,
    striker_id UUID REFERENCES team_players(id) ON DELETE SET NULL,
    non_striker_id UUID REFERENCES team_players(id) ON DELETE SET NULL,
    bowler_id UUID REFERENCES team_players(id) ON DELETE SET NULL,
    total_runs INTEGER DEFAULT 0,
    total_wickets INTEGER DEFAULT 0,
    legal_balls INTEGER DEFAULT 0,
    current_over INTEGER DEFAULT 0,
    current_ball INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'live',
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_innings_state_status ON innings_state(status);

COMMIT;
