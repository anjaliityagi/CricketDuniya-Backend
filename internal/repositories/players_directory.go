package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"fmt"
	"strings"
)

func GetPlayersDirectory(search string, limit int) ([]dto.PlayerDirectoryItem, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	query := `
    SELECT
    u.id,
    u.name AS name,
    u.phone_number,
    u.batting_style,
    u.bowling_style,
    COUNT(DISTINCT pms.match_id)::INT AS matches_played,
    COALESCE(SUM(pms.fantasy_points), 0)::INT AS points
   FROM users u
   LEFT JOIN team_players tp ON tp.player_id = u.id
   LEFT JOIN player_match_stats pms ON pms.team_player_id = tp.id`

	args := []interface{}{}
	if trimmedSearch := strings.TrimSpace(search); trimmedSearch != "" {
		query += ` WHERE u.name ILIKE $1 OR u.phone_number ILIKE $1`
		args = append(args, "%"+trimmedSearch+"%")
	}

	query += fmt.Sprintf(`GROUP BY
    u.id,
    u.name,
    u.phone_number,
    u.batting_style,
    u.bowling_style
    ORDER BY
    COALESCE(SUM(pms.fantasy_points), 0) DESC,
    COUNT(DISTINCT pms.match_id) DESC,
    u.name ASC
    LIMIT $%d`, len(args)+1)

	args = append(args, limit)
	var players []dto.PlayerDirectoryItem
	err := database.DB.Select(&players, query, args...)
	if err != nil {
		return nil, err
	}
	if players == nil {
		players = []dto.PlayerDirectoryItem{}
	}
	return players, nil
}
