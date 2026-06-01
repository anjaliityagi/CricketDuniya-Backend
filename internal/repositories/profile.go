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
		batting_style,
		bowling_style,
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
		batting_style = COALESCE($3, batting_style),
		bowling_style = COALESCE($4, bowling_style),
		updated_at = NOW()
	WHERE id = $1
	RETURNING
		id,
		name,
		phone_number,
		batting_style,
		bowling_style,
		created_at,
		updated_at
	`

	var user dto.UserProfileUser
	err := database.DB.QueryRowx(
		query,
		userID,
		req.Name,
		req.BattingStyle,
		req.BowlingStyle,
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
			m.winner_match_team_id,
			m.player_of_match_user_id,
			tp.team_id
		FROM player_match_stats pms
		LEFT JOIN matches m ON m.id = pms.match_id
		LEFT JOIN team_players tp ON tp.id = pms.team_player_id
		WHERE tp.player_id = $1
	)
	SELECT
		COUNT(DISTINCT match_id)::INT AS matches_played,
		COUNT(*) FILTER (WHERE winner_match_team_id = team_id)::INT AS won,
		COUNT(*) FILTER (WHERE winner_match_team_id IS NOT NULL AND winner_match_team_id <> team_id)::INT AS lost,
		COUNT(*) FILTER (WHERE player_of_match_user_id = $1)::INT AS mvps,
	COALESCE(
    ROUND(
        COUNT(*) FILTER (WHERE winner_match_team_id = team_id)::NUMERIC * 100 /
        NULLIF(COUNT(*) FILTER (WHERE winner_match_team_id IS NOT NULL), 0),
        2
    ),
    0
)::FLOAT AS win_percentage,
		COALESCE(SUM(fantasy_points), 0)::INT AS points
	FROM user_base
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
	SELECT COALESCE(
    ROUND(
        SUM(runs_scored)::NUMERIC /
        NULLIF(COUNT(CASE WHEN is_out THEN 1 END), 0),
        2
    ),
    ROUND(COALESCE(SUM(runs_scored), 0)::NUMERIC, 2)
)::FLOAT AS average,
		COALESCE(
    ROUND(
        SUM(runs_scored)::NUMERIC * 100 /
        NULLIF(SUM(balls_faced), 0),
        2
    ),
    0
)::FLOAT AS strike_rate,
		COALESCE(MAX(runs_scored),0)::INT AS high_score,
		COALESCE(SUM(runs_scored),0)::INT AS runs,
		COUNT(*)::INT AS innings,
		COUNT(*) FILTER (
			WHERE COALESCE(runs_scored, 0) = 0
			AND (COALESCE(balls_faced, 0) > 0 OR COALESCE(is_out, FALSE))
		)::INT AS ducks,
		COUNT(*) FILTER (WHERE COALESCE(runs_scored, 0) >= 30 AND COALESCE(runs_scored, 0) < 50)::INT AS thirties,
		COUNT(*) FILTER (WHERE COALESCE(runs_scored, 0) >= 50 AND COALESCE(runs_scored, 0) < 100)::INT AS fifties,
		COUNT(*) FILTER (WHERE COALESCE(runs_scored, 0) >= 100)::INT AS hundreds,
		COALESCE(SUM(fours),0)::INT AS fours,
		COALESCE(SUM(sixes),0)::INT AS sixes
	FROM player_match_stats pms
	LEFT JOIN team_players tp ON tp.id = pms.team_player_id
	WHERE tp.player_id = $1
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
        NULLIF(SUM(legal_balls_bowled), 0),
        2
    ),
    0
)::FLOAT AS economy
	FROM player_match_stats pms
	LEFT JOIN team_players tp ON tp.id = pms.team_player_id
	WHERE tp.player_id = $1
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
	FROM player_match_stats pms
	LEFT JOIN team_players tp ON tp.id = pms.team_player_id
	WHERE tp.player_id = $1
	`

	var stats dto.UserFieldingStats
	err := database.DB.Get(&stats, query, userID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func GetUserRecentMatches(userID string) ([]dto.UserRecentMatchPerformance, error) {
	query := `
	SELECT
		pms.match_id,
		m.match_date,
		CASE
			WHEN tp.team_id = m.team_a_id THEN COALESCE(ta.name, 'Team A')
			ELSE COALESCE(tb.name, 'Team B')
		END AS team_name,
		CASE
			WHEN tp.team_id = m.team_a_id THEN COALESCE(tb.name, 'Team B')
			ELSE COALESCE(ta.name, 'Team A')
		END AS opponent_name,
		CASE
			WHEN m.winner_match_team_id IS NULL THEN 'Pending'
			WHEN m.winner_match_team_id = tp.team_id THEN 'Won'
			ELSE 'Lost'
		END AS result,
		COALESCE(pms.runs_scored, 0) AS runs_scored,
		COALESCE(pms.balls_faced, 0) AS balls_faced,
		COALESCE(pms.wickets_taken, 0) AS wickets_taken,
		COALESCE(pms.runs_conceded, 0) AS runs_conceded,
		(
			FLOOR(COALESCE(pms.legal_balls_bowled, 0) / 6.0)
			+ MOD(COALESCE(pms.legal_balls_bowled, 0), 6)::NUMERIC / 10.0
		)::FLOAT AS overs_bowled,
		COALESCE(pms.catches, 0) AS catches,
		COALESCE(pms.stumping, 0) AS stumping,
		COALESCE(pms.runouts, 0) AS run_outs,
		COALESCE(pms.fantasy_points, 0) AS fantasy_points
	FROM player_match_stats pms
	LEFT JOIN matches m ON m.id = pms.match_id
	LEFT JOIN team_players tp ON tp.id = pms.team_player_id
	LEFT JOIN teams ta ON ta.id = m.team_a_id
	LEFT JOIN teams tb ON tb.id = m.team_b_id
	WHERE tp.player_id = $1
	ORDER BY COALESCE(m.completed_at, m.match_date, pms.updated_at) DESC, pms.updated_at DESC
	LIMIT 5
	`

	var matches []dto.UserRecentMatchPerformance
	err := database.DB.Select(&matches, query, userID)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
