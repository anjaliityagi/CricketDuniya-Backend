package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
)

func CalculateWinProbability(matchID string) (*dto.WinProbabilityResponse, error) {

	state, err := repositories.GetCurrentInningsState(matchID)
	if err != nil {
		return nil, err
	}

	totalBalls := state.OversPerInnings * 6
	ballsRemaining := totalBalls - state.LegalBalls
	wicketsRemaining := 10 - state.TotalWickets

	var runsRequired int
	if state.TargetRuns != nil {
		runsRequired = *state.TargetRuns - state.TotalRuns
	} else {
		runsRequired = 0
	}

	if runsRequired <= 0 {
		return &dto.WinProbabilityResponse{
			MatchID:                matchID,
			BattingTeamProbability: 100,
			BowlingTeamProbability: 0,
			Innings:                state.InningsNumber,
			CalculatedFrom: dto.ProbabilityFactors{
				RunsRequired:     runsRequired,
				BallsRemaining:   ballsRemaining,
				WicketsRemaining: wicketsRemaining,
			},
		}, nil
	}

	var probability float64

	if state.InningsNumber == 2 {

		probability = 50 +
			float64(wicketsRemaining*5) -
			float64(runsRequired)/2

		if probability > 100 {
			probability = 100
		}
		if probability < 0 {
			probability = 0
		}

	} else {

		matchScore := state.TotalRuns + (wicketsRemaining * 10)

		switch {
		case matchScore < 80:
			probability = 20
		case matchScore < 120:
			probability = 40
		case matchScore < 160:
			probability = 60
		case matchScore < 200:
			probability = 80
		default:
			probability = 90
		}
	}

	return &dto.WinProbabilityResponse{
		MatchID:                matchID,
		BattingTeamProbability: probability,
		BowlingTeamProbability: 100 - probability,
		Innings:                state.InningsNumber,
		CalculatedFrom: dto.ProbabilityFactors{
			RunsRequired:     runsRequired,
			BallsRemaining:   ballsRemaining,
			WicketsRemaining: wicketsRemaining,
		},
	}, nil
}
