package matchscoring

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
	"strings"
)

type Update struct {
	InningsRuns    int `json:"innings_runs"`
	InningsWickets int `json:"innings_wickets"`
}

func Process(req dto.BallRequest) (*Update, error) {
	tx, err := repositories.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	wicketInc := 0
	if req.IsWicket {
		wicketInc = 1
	}

	inningsRuns, inningsWickets, err := repositories.UpdateInningsTotalsTx(tx, req.InningsID.String(), req.TotalRuns, wicketInc)
	if err != nil {
		return nil, err
	}

	legalBall := !isExtraBall(req.BallType)

	if req.StrikerID.String() != "" {
		if err := repositories.UpsertBattingStatsTx(tx, req.MatchID.String(), req.StrikerID.String(), req.RunsOffBat, legalBall, req.IsBoundaryFour, req.IsBoundarySix, req.IsWicket); err != nil {
			return nil, err
		}
	}

	if req.BowlerID.String() != "" {
		conceded := req.TotalRuns - req.Byes - req.LegByes
		if conceded < 0 {
			conceded = 0
		}
		if err := repositories.UpsertBowlingStatsTx(tx, req.MatchID.String(), req.BowlerID.String(), conceded, legalBall, req.IsWicket); err != nil {
			return nil, err
		}
	}

	if req.FielderID != nil && req.FielderID.String() != "" && req.IsWicket {
		dismissalType := strings.ToLower(strings.TrimSpace(req.DismissalType))
		if err := repositories.UpsertFieldingStatsTx(
			tx,
			req.MatchID.String(),
			req.FielderID.String(),
			dismissalType == "caught",
			dismissalType == "stumped",
			dismissalType == "run_out",
		); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Update{
		InningsRuns:    inningsRuns,
		InningsWickets: inningsWickets,
	}, nil
}

func isExtraBall(ballType string) bool {
	bt := strings.ToLower(strings.TrimSpace(ballType))
	return bt == "wide" || bt == "no_ball" || bt == "dead_ball" || bt == "retired_hurt"
}
