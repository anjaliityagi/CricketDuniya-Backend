package scoring

import "CricketDuniya-Backend/internal/repositories"

func ApplyResultPoints(matchID string, winnerMatchTeamID string) error {
	tx, err := repositories.BeginTx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := repositories.ApplyResultPointsTx(tx, matchID, winnerMatchTeamID); err != nil {
		return err
	}

	return tx.Commit()
}

func ApplyNotOutBonus(matchID string, inningsID string) error {
	return repositories.ApplyNotOutBonus(matchID, inningsID)
}
