package scoring

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"

	"github.com/jmoiron/sqlx"
)

func PersistBallFantasy(req dto.BallRequest, ballEventID string, battingPoints, bowlingPoints, fieldingPoints int) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if req.StrikerID.String() != "" && battingPoints != 0 {
		if err := upsertFantasyPoints(tx, req.MatchID.String(), req.StrikerID.String(), battingPoints, "batting_points"); err != nil {
			return err
		}
		if err := insertPointEvent(tx, req.MatchID.String(), req.StrikerID.String(), ballEventID, "batting", "ball_batting_points", battingPoints); err != nil {
			return err
		}
	}

	if req.BowlerID.String() != "" && bowlingPoints != 0 {
		if err := upsertFantasyPoints(tx, req.MatchID.String(), req.BowlerID.String(), bowlingPoints, "bowling_points"); err != nil {
			return err
		}
		if err := insertPointEvent(tx, req.MatchID.String(), req.BowlerID.String(), ballEventID, "bowling", "ball_bowling_points", bowlingPoints); err != nil {
			return err
		}
	}

	if req.FielderID != nil && req.FielderID.String() != "" && fieldingPoints != 0 {
		if err := upsertFantasyPoints(tx, req.MatchID.String(), req.FielderID.String(), fieldingPoints, "fielding_points"); err != nil {
			return err
		}
		if err := insertPointEvent(tx, req.MatchID.String(), req.FielderID.String(), ballEventID, "fielding", "ball_fielding_points", fieldingPoints); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func upsertFantasyPoints(tx *sqlx.Tx, matchID, matchPlayerID string, points int, bucket string) error {
	if bucket != "batting_points" && bucket != "bowling_points" && bucket != "fielding_points" {
		return nil
	}
	query := `UPDATE player_match_stats
		SET` + bucket + ` = COALESCE(` + bucket + `, 0) + $3,
			fantasy_points = COALESCE(fantasy_points, 0) + $3,
			updated_at = NOW()
		WHERE match_id = $1
		  AND team_player_id = $2`
	res, err := tx.Exec(query, matchID, matchPlayerID, points)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO player_match_stats (match_id,team_player_id, ` + bucket + `, fantasy_points, updated_at)
		VALUES ($1, $2, $3, $3, NOW())
	`
	_, err = tx.Exec(insertQuery, matchID, matchPlayerID, points)
	return err
}

func insertPointEvent(tx *sqlx.Tx, matchID, matchPlayerID, ballEventID, category, ruleName string, points int) error {
	_, err := tx.Exec(`INSERT INTO point_events (
			match_id,
			user_id,
			ball_event_id,
			category,
			rule_name,
			points)
		SELECT
			$1,
			tp.player_id,
			$2,
			$3::point_category,
			$4,
			$5
		FROM team_players tp
		WHERE tp.id = $6
		  AND tp.player_id IS NOT NULL`,
		matchID, ballEventID, category, ruleName, points, matchPlayerID)
	return err
}
