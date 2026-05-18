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
		profile_image,
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
		profile_image = COALESCE($3, profile_image),
		batting_style = COALESCE($4, batting_style),
		bowling_style = COALESCE($5, bowling_style),
		updated_at = NOW()
	WHERE id = $1
	RETURNING
		id,
		name,
		phone_number,
		profile_image,
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
		req.ProfileImage,
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
	WITH user_matches AS (
		SELECT DISTINCT
			mp.match_id,
			mp.team_id,
			m.winner_team_id
		FROM match_players mp
		JOIN matches m ON m.id = mp.match_id
		WHERE mp.user_id = $1
	),
	user_points AS (
		SELECT
			COALESCE(SUM(pms.fantasy_points), 0) AS points
		FROM player_match_stats pms
		JOIN match_players mp ON mp.id = pms.player_id
		WHERE mp.user_id = $1
	),
	user_mvps AS (
		SELECT COUNT(*) AS mvps
		FROM player_match_stats pms
		JOIN match_players mp ON mp.id = pms.player_id
		WHERE mp.user_id = $1
		AND pms.fantasy_points = (
			SELECT MAX(match_stats.fantasy_points)
			FROM player_match_stats match_stats
			WHERE match_stats.match_id = pms.match_id
		)
	)
	SELECT
		COUNT(um.match_id)::INT AS matches_played,
		COUNT(CASE WHEN um.winner_team_id IS NOT NULL AND um.winner_team_id = um.team_id THEN 1 END)::INT AS won,
		COUNT(CASE WHEN um.winner_team_id IS NOT NULL AND um.winner_team_id <> um.team_id THEN 1 END)::INT AS lost,
		COALESCE((SELECT mvps FROM user_mvps), 0)::INT AS mvps,
		COALESCE(
			ROUND(
				COUNT(CASE WHEN um.winner_team_id IS NOT NULL AND um.winner_team_id = um.team_id THEN 1 END)::NUMERIC
				* 100 / NULLIF(COUNT(CASE WHEN um.winner_team_id IS NOT NULL THEN 1 END), 0),
				2
			),
			0
		)::FLOAT AS win_percentage,
		COALESCE((SELECT points FROM user_points), 0)::INT AS points
	FROM user_matches um
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
		COALESCE(ROUND(SUM(pms.runs)::NUMERIC / NULLIF(COUNT(CASE WHEN pms.is_out THEN 1 END), 0), 2), 0)::FLOAT AS average,
		COALESCE(ROUND(SUM(pms.runs)::NUMERIC * 100 / NULLIF(SUM(pms.balls_faced), 0), 2), 0)::FLOAT AS strike_rate,
		COALESCE(MAX(pms.runs), 0)::INT AS high_score,
		COALESCE(SUM(pms.runs), 0)::INT AS runs,
		COUNT(CASE WHEN pms.balls_faced > 0 OR pms.runs > 0 THEN 1 END)::INT AS innings,
		COALESCE(SUM(pms.fours), 0)::INT AS fours,
		COALESCE(SUM(pms.sixes), 0)::INT AS sixes
	FROM player_match_stats pms
	JOIN match_players mp ON mp.id = pms.player_id
	WHERE mp.user_id = $1
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
		COALESCE(SUM(pms.overs_bowled), 0)::FLOAT AS overs_bowled,
		COALESCE(SUM(pms.wickets), 0)::INT AS wickets,
		COALESCE(SUM(pms.runs_conceded), 0)::INT AS runs_conceded,
		COALESCE(SUM(pms.maidens), 0)::INT AS maidens,
		COALESCE(ROUND(SUM(pms.runs_conceded)::NUMERIC / NULLIF(SUM(pms.overs_bowled), 0), 2), 0)::FLOAT AS economy
	FROM player_match_stats pms
	JOIN match_players mp ON mp.id = pms.player_id
	WHERE mp.user_id = $1
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
		COALESCE(SUM(pms.catches), 0)::INT AS catches,
		COALESCE(SUM(pms.stumping), 0)::INT AS stumping,
		COALESCE(SUM(pms.run_outs), 0)::INT AS run_outs
	FROM player_match_stats pms
	JOIN match_players mp ON mp.id = pms.player_id
	WHERE mp.user_id = $1
	`

	var stats dto.UserFieldingStats
	err := database.DB.Get(&stats, query, userID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
