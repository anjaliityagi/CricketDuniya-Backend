package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/models"
)

type WinProbabilityState struct {
	MatchID         string
	InningsID       string
	InningsNumber   int
	TotalRuns       int
	TotalWickets    int
	LegalBalls      int
	OversPerInnings int
	TargetRuns      *int
}

func GetCurrentInningsState(matchID string) (*WinProbabilityState, error) {
	var innings models.Innings

	err := database.DB.Get(&innings, `
		SELECT id,match_id,innings_no,total_runs,total_wickets,target_runs
		FROM innings
		WHERE match_id = $1
		AND innings.completed_at IS NULL
		ORDER BY innings_no DESC
		LIMIT 1
	`, matchID)

	if err != nil {
		return nil, err
	}

	var state struct {
		LegalBalls int `db:"legal_balls"`
	}

	err = database.DB.Get(&state, `
		SELECT legal_balls
		FROM innings_state
		WHERE innings_id = $1
	`, innings.ID)

	if err != nil {
		return nil, err
	}

	var overs int
	_ = database.DB.Get(&overs, `
		SELECT overs_per_innings FROM matches WHERE id=$1
	`, matchID)

	return &WinProbabilityState{
		MatchID:         matchID,
		InningsID:       innings.ID,
		InningsNumber:   innings.InningsNumber,
		TotalRuns:       innings.TotalRuns,
		TotalWickets:    innings.Wickets,
		LegalBalls:      state.LegalBalls,
		OversPerInnings: overs,
		TargetRuns:      innings.TargetRuns,
	}, nil
}
