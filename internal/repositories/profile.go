package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
)

func GetUserProfileUser(userID string) (*dto.UserProfileUser, error) {
	query := `
	SELECT
		id,
		name,
		phone_number,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`

	var user dto.UserProfileUser
	err := database.DB.Get(&user, query, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserProfile(userID string, req dto.UpdateProfileRequest) (*dto.UserProfileUser, error) {
	query := `
	UPDATE users
	SET
		name = COALESCE($2, name),
		updated_at = NOW()
	WHERE id = $1
	RETURNING
		id,
		name,
		phone_number,
		created_at,
		updated_at
	`

	var user dto.UserProfileUser
	err := database.DB.QueryRowx(
		query,
		userID,
		req.Name,
	).StructScan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func GetUserProfileSummary(userID string) (*dto.UserProfileSummary, error) {

	query := `
	WITH user_base AS (
		SELECT 
			pms.match_id,
			pms.fantasy_points,
			pms.is_out,
		
			m.winner_match_team_id
		FROM player_match_stats pms
		LEFT JOIN matches m ON m.id = pms.match_id
		WHERE pms.user_id = $1
	),

	user_stats AS (
		SELECT
			COUNT(DISTINCT match_id) AS matches_played,
			COALESCE(SUM(fantasy_points), 0) AS total_points,
			COUNT(*) FILTER (WHERE is_out = true) AS dismissals
		FROM user_base
	),

	win_stats AS (
		SELECT
			COUNT(*) FILTER (WHERE winner_match_team_id IS NOT NULL) AS total_matches
-- 			COUNT(*) FILTER (WHERE winner_match_team_id = match_team_id) AS wins
		FROM user_base
	)

	SELECT
		COALESCE(us.matches_played, 0)::INT AS matches_played,
-- 		COALESCE(ws.wins, 0)::INT AS won,
-- 		COALESCE(ws.total_matches - ws.wins, 0)::INT AS lost,
		COALESCE(us.total_points, 0)::INT AS points,

		0::NUMERIC AS win_percentage

	FROM user_stats us
	CROSS JOIN win_stats ws
	`

	var stats dto.UserProfileSummary
	err := database.DB.Get(&stats, query, userID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
func GetUserBattingStats(userID string) (*dto.UserBattingStats, error) {

	query := `
	SELECT
		COALESCE(
			ROUND(
				SUM(runs_scored)::NUMERIC /
				NULLIF(COUNT(CASE WHEN is_out THEN 1 END), 0),
			2),
		0)::FLOAT AS average,

		COALESCE(
			ROUND(
				SUM(runs_scored)::NUMERIC * 100 /
				NULLIF(SUM(balls_faced), 0),
			2),
		0)::FLOAT AS strike_rate,

		COALESCE(MAX(runs_scored),0)::INT AS high_score,
		COALESCE(SUM(runs_scored),0)::INT AS runs,
		COUNT(*)::INT AS innings,
		COALESCE(SUM(fours),0)::INT AS fours,
		COALESCE(SUM(sixes),0)::INT AS sixes

	FROM player_match_stats
	WHERE user_id = $1
	`

	var stats dto.UserBattingStats
	err := database.DB.Get(&stats, query, userID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
func GetUserBowlingStats(userID string) (*dto.UserBowlingStats, error) {

	query := `
	SELECT
		(
			FLOOR(COALESCE(SUM(legal_balls_bowled), 0) / 6.0)
			+ MOD(COALESCE(SUM(legal_balls_bowled), 0), 6)::NUMERIC / 10.0
		)::FLOAT AS overs_bowled,
		COALESCE(SUM(wickets_taken),0)::INT AS wickets,
		COALESCE(SUM(runs_conceded),0)::INT AS runs_conceded,
		COALESCE(SUM(maidens),0)::INT AS maidens,

		COALESCE(
			ROUND(
				SUM(runs_conceded)::NUMERIC * 6 /
				NULLIF(SUM(legal_balls_bowled),0),
			2),
		0)::FLOAT AS economy

	FROM player_match_stats
	WHERE user_id = $1
	`

	var stats dto.UserBowlingStats
	err := database.DB.Get(&stats, query, userID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
func GetUserFieldingStats(userID string) (*dto.UserFieldingStats, error) {

	query := `
	SELECT
		COALESCE(SUM(catches),0)::INT AS catches,
		COALESCE(SUM(stumping),0)::INT AS stumping,
		COALESCE(SUM(runouts),0)::INT AS run_outs
	FROM player_match_stats
	WHERE user_id = $1
	`

	var stats dto.UserFieldingStats
	err := database.DB.Get(&stats, query, userID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
