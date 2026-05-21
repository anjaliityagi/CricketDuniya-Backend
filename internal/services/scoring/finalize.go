package scoring

import "CricketDuniya-Backend/internal/database"

func ApplyResultPoints(matchID string, winnerMatchTeamID string) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE player_match_stats pms
		SET
			result_points = COALESCE(result_points, 0) + 5,
			fantasy_points = COALESCE(fantasy_points, 0) + 5,
			updated_at = NOW()
		FROM match_team_players mtp
		WHERE pms.match_id = $1
		  AND pms.match_team_player_id = mtp.id
		  AND mtp.match_team_id = $2
		  AND mtp.is_playing_xi = TRUE
		  AND mtp.deleted_at IS NULL
	`, matchID, winnerMatchTeamID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE player_match_stats pms
		SET
			result_points = COALESCE(result_points, 0) - 5,
			fantasy_points = COALESCE(fantasy_points, 0) - 5,
			updated_at = NOW()
		FROM match_team_players mtp
		WHERE pms.match_id = $1
		  AND pms.match_team_player_id = mtp.id
		  AND mtp.match_team_id <> $2
		  AND mtp.is_playing_xi = TRUE
		  AND mtp.deleted_at IS NULL
	`, matchID, winnerMatchTeamID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func ApplyNotOutBonus(matchID string, inningsID string) error {
	_, err := database.DB.Exec(`
		UPDATE player_match_stats pms
		SET
			batting_points = COALESCE(batting_points, 0) + 5,
			fantasy_points = COALESCE(fantasy_points, 0) + 5,
			updated_at = NOW()
		WHERE pms.match_id = $1
		  AND pms.match_team_player_id IN (
			  SELECT DISTINCT be.striker_match_player_id
			  FROM ball_events be
			  WHERE be.match_id = $1
			    AND be.innings_id = $2
			    AND be.striker_match_player_id IS NOT NULL
			    AND be.is_deleted = FALSE
		  )
		  AND COALESCE(pms.is_out, FALSE) = FALSE
	`, matchID, inningsID)
	return err
}
