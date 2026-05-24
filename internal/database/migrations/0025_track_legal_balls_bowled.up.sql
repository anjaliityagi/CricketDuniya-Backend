BEGIN;

ALTER TABLE player_match_stats
    ADD COLUMN IF NOT EXISTS legal_balls_bowled INTEGER NOT NULL DEFAULT 0;

WITH bowling_stats AS (
    SELECT
        be.match_id,
        be.bowler_id AS team_player_id,
        COUNT(*) FILTER (
            WHERE COALESCE(be.is_deleted, FALSE) = FALSE
              AND COALESCE(be.ball_type, 'normal') NOT IN ('wide', 'no_ball')
        )::INTEGER AS legal_balls_bowled
    FROM ball_events be
    WHERE be.bowler_id IS NOT NULL
    GROUP BY be.match_id, be.bowler_id
)
UPDATE player_match_stats pms
SET legal_balls_bowled = bs.legal_balls_bowled,
    overs_bowled = FLOOR(bs.legal_balls_bowled / 6.0)
        + MOD(bs.legal_balls_bowled, 6)::NUMERIC / 10.0
FROM bowling_stats bs
WHERE pms.match_id = bs.match_id
  AND pms.team_player_id = bs.team_player_id;

UPDATE player_match_stats
SET overs_bowled = FLOOR(COALESCE(legal_balls_bowled, 0) / 6.0)
    + MOD(COALESCE(legal_balls_bowled, 0), 6)::NUMERIC / 10.0
WHERE overs_bowled IS DISTINCT FROM (
    FLOOR(COALESCE(legal_balls_bowled, 0) / 6.0)
    + MOD(COALESCE(legal_balls_bowled, 0), 6)::NUMERIC / 10.0
);

COMMIT;
