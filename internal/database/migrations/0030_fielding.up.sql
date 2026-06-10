BEGIN;

UPDATE player_match_stats pms
SET catches = 0,
    stumping = 0,
    runouts = 0,
    updated_at = NOW()
WHERE EXISTS (
    SELECT 1
    FROM ball_events be
    WHERE be.match_id = pms.match_id
      AND be.is_deleted = FALSE
);

INSERT INTO player_match_stats (match_id, player_id, team_player_id, catches, stumping, runouts, updated_at)
SELECT
    be.match_id,
    tp.player_id,be.fielder_id,
    SUM(CASE WHEN be.dismissal_type = 'caught' THEN 1 ELSE 0 END)::INT,
    SUM(CASE WHEN be.dismissal_type = 'stumped' THEN 1 ELSE 0 END)::INT,
    SUM(CASE WHEN be.dismissal_type = 'run_out' THEN 1 ELSE 0 END)::INT,
    NOW()
FROM ball_events be
JOIN team_players tp ON tp.id = be.fielder_id
WHERE be.is_deleted = FALSE
  AND be.is_wicket = TRUE
  AND be.fielder_id IS NOT NULL
  AND be.dismissal_type IN ('caught', 'stumped', 'run_out')
  AND tp.player_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1
      FROM player_match_stats pms
      WHERE pms.match_id = be.match_id
        AND pms.team_player_id = be.fielder_id
  )
GROUP BY be.match_id, tp.player_id, be.fielder_id;

UPDATE player_match_stats pms
SET catches = stats.catches,
    stumping = stats.stumping,
    runouts = stats.runouts,
    updated_at = NOW()
FROM (
    SELECT
        match_id,
        fielder_id,
        SUM(CASE WHEN dismissal_type = 'caught' THEN 1 ELSE 0 END)::INT AS catches,
        SUM(CASE WHEN dismissal_type = 'stumped' THEN 1 ELSE 0 END)::INT AS stumping,
        SUM(CASE WHEN dismissal_type = 'run_out' THEN 1 ELSE 0 END)::INT AS runouts
    FROM ball_events
    WHERE is_deleted = FALSE
      AND is_wicket = TRUE
      AND fielder_id IS NOT NULL
      AND dismissal_type IN ('caught', 'stumped', 'run_out')
    GROUP BY match_id, fielder_id
) stats
WHERE pms.match_id = stats.match_id
  AND pms.team_player_id = stats.fielder_id;

COMMIT;

