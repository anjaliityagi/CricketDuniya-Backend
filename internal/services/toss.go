package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
	"errors"
)

func PerformToss(req dto.TossRequest) ([]dto.InningsResponse, error) {

	matchTeams, err := repositories.GetMatchTeams(req.MatchID)

	if err != nil {
		return nil, err
	}

	if len(matchTeams) != 2 {
		return nil, errors.New("match teams are not ready")
	}

	if req.Decision != "bat" && req.Decision != "bowl" {
		return nil, errors.New("invalid toss decision")
	}

	tossWinnerInputID := req.TossWinnerTeamID
	//if tossWinnerInputID == "" {
	//	tossWinnerInputID = req.TossWinnerTeamID
	//}
	if tossWinnerInputID == "" {
		return nil, errors.New("toss winner team id is required")
	}

	resolvedWinnerTeamID, err := repositories.ResolveMatchTeamID(req.MatchID, tossWinnerInputID)
	if err != nil {
		return nil, errors.New("toss winner team is not part of this match")
	}

	var winnerIdx = -1
	for i := range matchTeams {
		if matchTeams[i].ID == resolvedWinnerTeamID {
			winnerIdx = i
			break
		}
	}
	if winnerIdx == -1 {
		return nil, errors.New("toss winner team is not part of this match")
	}

	loserIdx := 1 - winnerIdx
	battingTeamID := matchTeams[winnerIdx].ID
	bowlingTeamID := matchTeams[loserIdx].ID
	if req.Decision == "bowl" {
		battingTeamID, bowlingTeamID = bowlingTeamID, battingTeamID
	}

	err = repositories.UpdateToss(
		req.MatchID,
		resolvedWinnerTeamID,
		req.Decision,
	)

	if err != nil {
		return nil, err
	}

	err = repositories.CreateInnings(req.MatchID, battingTeamID, bowlingTeamID, 1)

	if err != nil {
		return nil, err
	}

	innings, err := repositories.GetMatchInnings(req.MatchID)
	if err != nil {
		return nil, err
	}

	return innings, nil
}
