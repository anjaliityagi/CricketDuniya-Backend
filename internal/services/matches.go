package services

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
	"CricketDuniya-Backend/internal/services/scoring"
	"errors"
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

func CreateMatchWithTeams(req dto.CreateMatchWithTeamsRequest, hostUserID string) (*models.Team, *models.Team, *models.Match, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, nil, nil, err
	}
	defer tx.Rollback()

	teamA := &models.Team{
		Name:      &req.TeamAName,
		CreatedBy: &hostUserID,
	}
	if err := repositories.CreateTeamTx(tx, teamA); err != nil {
		return nil, nil, nil, err
	}

	teamB := &models.Team{
		Name:      &req.TeamBName,
		CreatedBy: &hostUserID,
	}
	if err := repositories.CreateTeamTx(tx, teamB); err != nil {
		return nil, nil, nil, err
	}

	var parsedDate *time.Time
	if req.MatchDate != "" {
		t, err := time.Parse(time.RFC3339, req.MatchDate)
		if err != nil {
			return nil, nil, nil, err
		}
		parsedDate = &t
	}

	match := &models.Match{
		TeamAID:         teamA.ID,
		TeamBID:         teamB.ID,
		HostUserID:      hostUserID,
		Location:        &req.Location,
		MatchDate:       parsedDate,
		OversPerInnings: req.OversPerInnings,
		Status:          "scheduled",
	}
	if err := repositories.CreateMatchTx(tx, match); err != nil {
		return nil, nil, nil, err
	}

	if _, _, err = repositories.CreateMatchSnapshots(match.ID, *teamA.ID, *teamB.ID); err != nil {
		return nil, nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, err
	}

	return teamA, teamB, match, nil
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

func SetFirstPickTeam(matchID string, firstPickTeamID string) (*dto.MatchResponse, error) {
	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		return nil, err
	}

	if match.TeamAID == nil || match.TeamBID == nil {
		return nil, errors.New("match teams are not ready")
	}

	if firstPickTeamID != *match.TeamAID && firstPickTeamID != *match.TeamBID {
		return nil, errors.New("first pick team must be one of the match teams")
	}

	if err := repositories.UpdateFirstPickTeam(matchID, firstPickTeamID); err != nil {
		return nil, err
	}

	return repositories.GetMatchDetailByID(matchID)
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
