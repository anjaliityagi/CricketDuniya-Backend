package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
	"errors"
)

func PerformToss(req dto.TossRequest) error {

	match, err := repositories.GetMatchByID(req.MatchID)
	if err != nil {
		return err
	}

	if match.TeamAID == nil || match.TeamBID == nil {
		return errors.New("teams are not assigned to match")
	}

	var battingTeamID string
	var bowlingTeamID string

	if req.Decision == "bat" {

		battingTeamID = req.TossWinnerTeamID

		if req.TossWinnerTeamID == *match.TeamAID {

			bowlingTeamID = *match.TeamBID

		} else {

			bowlingTeamID = *match.TeamAID
		}

	} else if req.Decision == "bowl" {

		bowlingTeamID = req.TossWinnerTeamID

		if req.TossWinnerTeamID == *match.TeamAID {

			battingTeamID = *match.TeamBID

		} else {

			battingTeamID = *match.TeamAID
		}

	} else {

		return errors.New("invalid toss decision")
	}

	err = repositories.UpdateToss(
		req.MatchID,
		req.TossWinnerTeamID,
		req.Decision,
	)

	if err != nil {
		return err
	}

	err = repositories.CreateInnings(
		req.MatchID,
		battingTeamID,
		bowlingTeamID,
		1,
	)

	if err != nil {
		return err
	}

	return nil
}
