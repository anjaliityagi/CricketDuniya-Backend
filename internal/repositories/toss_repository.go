package repositories

import (
	"CricketDuniya-Backend/internal/database"
)

type Match struct {
	ID      string
	TeamAID string
	TeamBID string
}

type MatchTeam struct {
	ID     string  `db:"id"`
	TeamID *string `db:"team_id"`
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

func GetMatchTeams(matchID string) ([]MatchTeam, error) {
	query := `SELECT id, team_id
		FROM match_teams
		WHERE match_id = $1
		  AND deleted_at IS NULL
		ORDER BY created_at ASC`

	var rows []MatchTeam
	if err := database.DB.Select(&rows, query, matchID); err != nil {
		return nil, err
	}

	return rows, nil
}

func UpdateToss(matchID string, tossWinnerMatchTeamID string, decision string) error {

	query := `UPDATE matches
	SET
		toss_winner_match_team_id = $1,
		toss_winner_team_id = (
			SELECT team_id
			FROM match_teams
			WHERE id = $1
			LIMIT 1
		),
		toss_decision = $2,
		status = 'live'
	WHERE id = $3`

	_, err := database.DB.Exec(
		query,
		tossWinnerMatchTeamID,
		decision,
		matchID,
	)

	return err
}

func CreateInnings(
	matchID string,
	battingMatchTeamID string,
	bowlingMatchTeamID string,
	inningsNo int,
) error {

	query := `
	INSERT INTO innings (
		match_id,
		batting_match_team_id,
		bowling_match_team_id,
		innings_no,
		started_at
	)
	VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := database.DB.Exec(
		query,
		matchID,
		battingMatchTeamID,
		bowlingMatchTeamID,
		inningsNo,
	)

	return err
}
