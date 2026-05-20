package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
)

func SaveBall(req dto.BallRequest) error {
	query := `
	INSERT INTO ball_events (
		id,
		innings_id,
		match_id,
		striker_id,
		non_striker_id,
		bowler_id,
		ball_no,
		delivery_no,
		ball_type,
		runs_scored,
		runs_off_bat,
		extras,
		total_runs,
		is_dot_ball,
		is_boundary_four,
		is_boundary_six,
		is_wicket,
		dismissed_player_id,
		fielder_id,
		wides,
		no_balls,
		byes,
		leg_byes,
		created_at
	)
	VALUES (
		gen_random_uuid(),
		$1,$2,$3,$4,$5,
		$6,$7,$8,$9,$10,$11,$12,
		$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,
		NOW()
	)
	`

	_, err := database.DB.Exec(query,
		req.InningsID,
		req.MatchID,
		req.StrikerID,
		req.NonStrikerID,
		req.BowlerID,
		req.BallNo,
		req.DeliveryNo,
		req.BallType,
		req.RunsScored,
		req.RunsOffBat,
		req.Extras,
		req.TotalRuns,
		req.IsDotBall,
		req.IsBoundaryFour,
		req.IsBoundarySix,
		req.IsWicket,
		req.DismissedPlayerID,
		req.FielderID,
		req.Wides,
		req.NoBalls,
		req.Byes,
		req.LegByes,
	)

	return err
}
