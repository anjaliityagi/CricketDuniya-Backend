package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
	"CricketDuniya-Backend/internal/services/scoring"
	"time"
)

func CreateMatch(req dto.CreateMatchRequest, hostUserID string) (*models.Match, error) {

	var parsedDate *time.Time

	if req.MatchDate != "" {

		t, err := time.Parse(time.RFC3339, req.MatchDate)
		if err != nil {
			return nil, err
		}

		parsedDate = &t
	}

	match := &models.Match{
		TeamAID:         &req.TeamAID,
		TeamBID:         &req.TeamBID,
		HostUserID:      hostUserID,
		Location:        &req.Location,
		MatchDate:       parsedDate,
		OversPerInnings: req.OversPerInnings,
		Status:          "scheduled",
	}

	err := repositories.CreateMatch(match)
	if err != nil {
		return nil, err
	}

	_, _, err = repositories.CreateMatchSnapshots(match.ID, req.TeamAID, req.TeamBID)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func GetAllMatches(query dto.GetMatchesQuery) ([]dto.MatchResponse, error) {

	return repositories.GetAllMatches(query)
}

func GetMatchByID(id string) (*models.Match, error) {
	return repositories.GetMatchByID(id)
}

func GetMatchDetailByID(id string) (*dto.MatchResponse, error) {
	return repositories.GetMatchDetailByID(id)
}

func GetMatchInnings(matchID string) ([]dto.InningsResponse, error) {
	return repositories.GetMatchInnings(matchID)
}

func GetMatchScorecard(matchID string) (*dto.MatchScorecardResponse, error) {
	return repositories.GetMatchScorecard(matchID)
}

func StartMatch(matchID string) ([]dto.InningsResponse, error) {
	return repositories.StartMatch(matchID)
}

func GetMatchSquad(matchID string) ([]dto.MatchSquadPlayer, error) {
	return repositories.GetMatchSquad(matchID)
}

func UpdateMatchLineup(matchID string, req dto.UpdateMatchLineupRequest) error {
	return repositories.UpdateMatchLineup(matchID, req.Players)
}

func CompleteMatch(matchID, winnerTeamOrMatchTeamID string) error {
	winnerMatchTeamID, err := repositories.ResolveMatchTeamID(matchID, winnerTeamOrMatchTeamID)
	if err != nil {
		return err
	}

	inningsIDs, err := repositories.GetMatchInningsIDs(matchID)
	if err != nil {
		return err
	}

	for _, inningsID := range inningsIDs {
		if err := scoring.ApplyNotOutBonus(matchID, inningsID); err != nil {
			return err
		}
	}

	if err := scoring.ApplyResultPoints(matchID, winnerMatchTeamID); err != nil {
		return err
	}

	return repositories.FinalizeMatch(matchID, winnerMatchTeamID)
}
