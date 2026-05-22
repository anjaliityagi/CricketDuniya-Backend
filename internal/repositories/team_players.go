package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
)

func GetPlayersByTeamID(teamID string) ([]dto.TeamPlayerResponse, error) {
	var players []dto.TeamPlayerResponse

	query := `
	SELECT
		tp.id,
		tp.team_id,
		tp.player_id,
		COALESCE(u.name, '') AS name,
		u.phone_number,
		COALESCE(tp.is_captain, FALSE) AS is_captain,
		COALESCE(tp.is_wicket_keeper, FALSE) AS is_wicket_keeper,
		COALESCE(tp.is_substitute, FALSE) AS is_substitute,
		tp.batting_order,
		tp.created_at,
		NULL::TIMESTAMP AS removed_at
	FROM team_players tp
	LEFT JOIN users u ON u.id = tp.player_id
	WHERE tp.team_id = $1
	  AND tp.deleted_at IS NULL
	ORDER BY tp.batting_order NULLS LAST, tp.created_at ASC
	`

	err := database.DB.Select(&players, query, teamID)
	if err != nil {
		return nil, err
	}

	return players, nil
}
