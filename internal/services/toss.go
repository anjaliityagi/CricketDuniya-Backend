package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
	"errors"
)

func PerformToss(req dto.TossRequest) error {

	matchTeams, err := repositories.GetMatchTeams(req.MatchID)

	if err != nil {
		return err
	}

	if len(matchTeams) != 2 {
		return errors.New("match teams are not ready")
	}

	if req.Decision != "bat" && req.Decision != "bowl" {
		return errors.New("invalid toss decision")
	}

	var winnerIdx = -1
	for i := range matchTeams {
		if matchTeams[i].ID == req.TossWinnerTeamID {
			winnerIdx = i
			break
		}
	}
	if winnerIdx == -1 {
		return errors.New("toss winner team is not part of this match")
	}

	loserIdx := 1 - winnerIdx
	battingTeamID := matchTeams[winnerIdx].ID
	bowlingTeamID := matchTeams[loserIdx].ID
	if req.Decision == "bowl" {
		battingTeamID, bowlingTeamID = bowlingTeamID, battingTeamID
	}

	err = repositories.UpdateToss(
		req.MatchID,
		req.TossWinnerTeamID,
		req.Decision,
	)

	if err != nil {
		return err
	}

	err = repositories.CreateInnings(req.MatchID, battingTeamID, bowlingTeamID, 1)

	if err != nil {
		return err
	}

	return nil
}
