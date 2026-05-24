BEGIN;

ALTER TABLE ball_events
    ADD COLUMN IF NOT EXISTS runs_off_bat INTEGER DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_boundary_four BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_boundary_six BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_wicket BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS dismissed_player_id UUID REFERENCES team_players(id),
    ADD COLUMN IF NOT EXISTS fielder_id UUID REFERENCES team_players(id),
    ADD COLUMN IF NOT EXISTS is_dot_ball BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS wides INTEGER DEFAULT 0,
    ADD COLUMN IF NOT EXISTS no_balls INTEGER DEFAULT 0,
    ADD COLUMN IF NOT EXISTS byes INTEGER DEFAULT 0,
    ADD COLUMN IF NOT EXISTS leg_byes INTEGER DEFAULT 0;

UPDATE ball_events
SET runs_off_bat = CASE
        WHEN COALESCE(ball_type, 'normal') IN ('wide', 'bye', 'leg_bye') THEN 0
        ELSE GREATEST(COALESCE(total_runs, 0) - COALESCE(extras, 0), 0)
    END
WHERE runs_off_bat IS NULL
   OR runs_off_bat = 0;

UPDATE ball_events
SET is_boundary_four = (COALESCE(runs_off_bat, 0) = 4)
WHERE is_boundary_four IS DISTINCT FROM (COALESCE(runs_off_bat, 0) = 4);

UPDATE ball_events
SET is_boundary_six = (COALESCE(runs_off_bat, 0) = 6)
WHERE is_boundary_six IS DISTINCT FROM (COALESCE(runs_off_bat, 0) = 6);

UPDATE ball_events
SET is_wicket = CASE
        WHEN dismissal_type IS NOT NULL THEN TRUE
        ELSE FALSE
    END
WHERE is_wicket IS DISTINCT FROM CASE
        WHEN dismissal_type IS NOT NULL THEN TRUE
        ELSE FALSE
    END;

UPDATE ball_events
SET is_dot_ball = (COALESCE(total_runs, 0) = 0)
WHERE is_dot_ball IS DISTINCT FROM (COALESCE(total_runs, 0) = 0);

UPDATE ball_events
SET wides = CASE WHEN COALESCE(ball_type, 'normal') = 'wide' THEN COALESCE(extras, 0) ELSE COALESCE(wides, 0) END
WHERE COALESCE(ball_type, 'normal') = 'wide';

UPDATE ball_events
SET no_balls = CASE WHEN COALESCE(ball_type, 'normal') = 'no_ball' THEN COALESCE(extras, 0) ELSE COALESCE(no_balls, 0) END
WHERE COALESCE(ball_type, 'normal') = 'no_ball';

UPDATE ball_events
SET byes = CASE WHEN COALESCE(ball_type, 'normal') = 'bye' THEN COALESCE(extras, 0) ELSE COALESCE(byes, 0) END
WHERE COALESCE(ball_type, 'normal') = 'bye';

UPDATE ball_events
SET leg_byes = CASE WHEN COALESCE(ball_type, 'normal') = 'leg_bye' THEN COALESCE(extras, 0) ELSE COALESCE(leg_byes, 0) END
WHERE COALESCE(ball_type, 'normal') = 'leg_bye';

ALTER TABLE ball_events
    ALTER COLUMN runs_off_bat SET DEFAULT 0,
    ALTER COLUMN is_boundary_four SET DEFAULT FALSE,
    ALTER COLUMN is_boundary_six SET DEFAULT FALSE,
    ALTER COLUMN is_wicket SET DEFAULT FALSE,
    ALTER COLUMN is_dot_ball SET DEFAULT FALSE,
    ALTER COLUMN wides SET DEFAULT 0,
    ALTER COLUMN no_balls SET DEFAULT 0,
    ALTER COLUMN byes SET DEFAULT 0,
    ALTER COLUMN leg_byes SET DEFAULT 0;

COMMIT;
