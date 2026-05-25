package scoring

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
)

func PersistBallFantasy(req dto.BallRequest, ballEventID string, battingPoints, bowlingPoints, fieldingPoints int) error {
	tx, err := repositories.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if req.StrikerID.String() != "" && battingPoints != 0 {
		if err := repositories.UpsertFantasyPointsTx(tx, req.MatchID.String(), req.StrikerID.String(), battingPoints, "batting_points"); err != nil {
			return err
		}
		if err := repositories.InsertPointEventTx(tx, req.MatchID.String(), req.StrikerID.String(), ballEventID, "batting", "ball_batting_points", battingPoints); err != nil {
			return err
		}
	}

	if req.BowlerID.String() != "" && bowlingPoints != 0 {
		if err := repositories.UpsertFantasyPointsTx(tx, req.MatchID.String(), req.BowlerID.String(), bowlingPoints, "bowling_points"); err != nil {
			return err
		}
		if err := repositories.InsertPointEventTx(tx, req.MatchID.String(), req.BowlerID.String(), ballEventID, "bowling", "ball_bowling_points", bowlingPoints); err != nil {
			return err
		}
	}

	if req.FielderID != nil && req.FielderID.String() != "" && fieldingPoints != 0 {
		if err := repositories.UpsertFantasyPointsTx(tx, req.MatchID.String(), req.FielderID.String(), fieldingPoints, "fielding_points"); err != nil {
			return err
		}
		if err := repositories.InsertPointEventTx(tx, req.MatchID.String(), req.FielderID.String(), ballEventID, "fielding", "ball_fielding_points", fieldingPoints); err != nil {
			return err
		}
	}

	return tx.Commit()
}
