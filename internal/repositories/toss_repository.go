package repositories

import (
	"CricketDuniya-Backend/internal/database"
)

type Match struct {
	ID      string
	TeamAID string
	TeamBID string
}

//func GetMatchByID(matchID string) (*Match, error) {
//
//	query := `
//	SELECT id, team_a_id, team_b_id
//	FROM matches
//	WHERE id = $1
//	`
//
//	row := database.DB.QueryRow(query, matchID)
//
//	var match Match
//
//	err := row.Scan(
//		&match.ID,
//		&match.TeamAID,
//		&match.TeamBID,
//	)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return &match, nil
//}

func UpdateToss(matchID string, tossWinnerTeamID string, decision string) error {

	query := `
	UPDATE matches
	SET
		toss_winner_team_id = $1,
		toss_decision = $2,
		status = 'live'
	WHERE id = $3
	`

	_, err := database.DB.Exec(
		query,
		tossWinnerTeamID,
		decision,
		matchID,
	)

	return err
}

func CreateInnings(
	matchID string,
	battingTeamID string,
	bowlingTeamID string,
	inningsNumber int,
) error {

	query := `
	INSERT INTO innings (
		match_id,
		batting_team_id,
		bowling_team_id,
		innings_number
	)
	VALUES ($1, $2, $3, $4)
	`

	_, err := database.DB.Exec(
		query,
		matchID,
		battingTeamID,
		bowlingTeamID,
		inningsNumber,
	)

	return err
}
