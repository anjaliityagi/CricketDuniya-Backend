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

func GetMatchTeams(matchID string) ([]MatchTeam, error) {
	query := `
		SELECT team_a_id AS id, team_a_id AS team_id
		FROM matches
		WHERE id = $1
		  AND team_a_id IS NOT NULL
		UNION ALL
		SELECT team_b_id AS id, team_b_id AS team_id
		FROM matches
		WHERE id = $1
		  AND team_b_id IS NOT NULL
	`

	var rows []MatchTeam
	if err := database.DB.Select(&rows, query, matchID); err != nil {
		return nil, err
	}

	return rows, nil
}

func UpdateToss(matchID string, tossWinnerTeamID string, decision string) error {

	query := `UPDATE matches
	SET

		toss_winner_team_id = $1,
		toss_decision = $2,
		status = 'live'
	WHERE id = $3`

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
	battingMatchTeamID string,
	bowlingMatchTeamID string,
	inningsNo int,
) error {

	query := `
	INSERT INTO innings (
		match_id,
		batting_team_id,
		bowling_team_id,
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
