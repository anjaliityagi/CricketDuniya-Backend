package scoring

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"strings"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

type scoreContext struct {
	prevConsecutiveBatterFours int
	prevConsecutiveBatterSixes int

	prevConsecutiveBowlerFours         int
	prevConsecutiveBowlerSixes         int
	prevConsecutiveBowlerExtras        int
	prevConsecutiveBowlerWicketsInOver int

	batterRunsBefore            int
	batterBallsBefore           int
	batterStrikeRateBonusBefore int

	bowlerLegalBallsBefore   int
	bowlerRunsConcededBefore int
	bowlerEconomyBonusBefore int

	overLegalBallsBefore   int
	overRunsConcededBefore int
	overBoundariesBefore   int
}

func (e *Engine) Process(req dto.BallRequest) (bat, bowl, field int, err error) {
	ctx, err := buildScoreContext(req)
	if err != nil {
		return 0, 0, 0, err
	}

	bat = batting(req, ctx)
	bowl = bowling(req, ctx)
	field = fielding(req)
	return bat, bowl, field, nil
}

func buildScoreContext(req dto.BallRequest) (scoreContext, error) {
	ctx := scoreContext{}

	// batter streaks
	if req.StrikerID.String() != "" {
		n, err := countConsecutiveBatterBoundary(req, "four")
		if err != nil {
			return ctx, err
		}
		ctx.prevConsecutiveBatterFours = n

		n, err = countConsecutiveBatterBoundary(req, "six")
		if err != nil {
			return ctx, err
		}
		ctx.prevConsecutiveBatterSixes = n

		if err := database.DB.QueryRowx(`
			SELECT
				COALESCE(SUM(CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END), 0) AS runs_before,
				COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball') THEN 1 ELSE 0 END), 0) AS balls_before
			FROM ball_events
			WHERE match_id = $1
			  AND innings_id = $2
			  AND striker_id = $3
			  AND is_deleted = FALSE
		`, req.MatchID, req.InningsID, req.StrikerID).Scan(&ctx.batterRunsBefore, &ctx.batterBallsBefore); err != nil {
			return ctx, err
		}
		ctx.batterStrikeRateBonusBefore = strikeRateBonus(ctx.batterRunsBefore, ctx.batterBallsBefore)
	}

	// bowler streaks + stats
	if req.BowlerID.String() != "" {
		n, err := countConsecutiveBowlerConcededBoundary(req, "four")
		if err != nil {
			return ctx, err
		}
		ctx.prevConsecutiveBowlerFours = n

		n, err = countConsecutiveBowlerConcededBoundary(req, "six")
		if err != nil {
			return ctx, err
		}
		ctx.prevConsecutiveBowlerSixes = n

		n, err = countConsecutiveBowlerExtras(req)
		if err != nil {
			return ctx, err
		}
		ctx.prevConsecutiveBowlerExtras = n

		n, err = countConsecutiveBowlerWicketsInOver(req)
		if err != nil {
			return ctx, err
		}
		ctx.prevConsecutiveBowlerWicketsInOver = n

		if err := database.DB.QueryRowx(`
			SELECT
				COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball') THEN 1 ELSE 0 END), 0) AS legal_balls_before,
				COALESCE(SUM(total_runs - byes - leg_byes), 0) AS runs_conceded_before
			FROM ball_events
			WHERE match_id = $1
			  AND innings_id = $2
			  AND bowler_id = $3
			  AND is_deleted = FALSE
		`, req.MatchID, req.InningsID, req.BowlerID).Scan(&ctx.bowlerLegalBallsBefore, &ctx.bowlerRunsConcededBefore); err != nil {
			return ctx, err
		}
		ctx.bowlerEconomyBonusBefore = economyBonus(ctx.bowlerRunsConcededBefore, ctx.bowlerLegalBallsBefore)

		if err := database.DB.QueryRowx(`
			SELECT
				COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball') THEN 1 ELSE 0 END), 0) AS over_legal_before,
				COALESCE(SUM(total_runs - byes - leg_byes), 0) AS over_runs_before,
				COALESCE(SUM(CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 WHEN GREATEST(total_runs - extras, 0) IN (4, 6) THEN 1 ELSE 0 END), 0) AS over_boundaries_before
			FROM ball_events
			WHERE match_id = $1
			  AND innings_id = $2
			  AND bowler_id = $3
			  AND ball_no = $4
			  AND is_deleted = FALSE
		`, req.MatchID, req.InningsID, req.BowlerID, req.BallNo).Scan(&ctx.overLegalBallsBefore, &ctx.overRunsConcededBefore, &ctx.overBoundariesBefore); err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func countConsecutiveBatterBoundary(req dto.BallRequest, boundaryType string) (int, error) {
	rows, err := database.DB.Queryx(`
		SELECT CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END AS runs_off_bat
		FROM ball_events
		WHERE match_id = $1
		  AND innings_id = $2
		  AND striker_id = $3
		  AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 24
	`, req.MatchID, req.InningsID, req.StrikerID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var runsOffBat int
		if err := rows.Scan(&runsOffBat); err != nil {
			return 0, err
		}
		if boundaryType == "four" && runsOffBat == 4 {
			count++
			continue
		}
		if boundaryType == "six" && runsOffBat == 6 {
			count++
			continue
		}
		break
	}
	return count, nil
}

func countConsecutiveBowlerConcededBoundary(req dto.BallRequest, boundaryType string) (int, error) {
	rows, err := database.DB.Queryx(`
		SELECT CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END AS runs_off_bat
		FROM ball_events
		WHERE match_id = $1
		  AND innings_id = $2
		  AND bowler_id = $3
		  AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 24
	`, req.MatchID, req.InningsID, req.BowlerID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var runsOffBat int
		if err := rows.Scan(&runsOffBat); err != nil {
			return 0, err
		}
		if boundaryType == "four" && runsOffBat == 4 {
			count++
			continue
		}
		if boundaryType == "six" && runsOffBat == 6 {
			count++
			continue
		}
		break
	}
	return count, nil
}

func countConsecutiveBowlerExtras(req dto.BallRequest) (int, error) {
	rows, err := database.DB.Queryx(`
		SELECT ball_type
		FROM ball_events
		WHERE match_id = $1
		  AND innings_id = $2
		  AND bowler_id = $3
		  AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 24
	`, req.MatchID, req.InningsID, req.BowlerID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return 0, err
		}
		bt := strings.ToLower(strings.TrimSpace(t))
		if bt == "wide" || bt == "no_ball" {
			count++
			continue
		}
		break
	}
	return count, nil
}

func countConsecutiveBowlerWicketsInOver(req dto.BallRequest) (int, error) {
	rows, err := database.DB.Queryx(`
		SELECT is_wicket
		FROM ball_events
		WHERE match_id = $1
		  AND innings_id = $2
		  AND bowler_id = $3
		  AND ball_no = $4
		  AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 24
	`, req.MatchID, req.InningsID, req.BowlerID, req.BallNo)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var isWicket bool
		if err := rows.Scan(&isWicket); err != nil {
			return 0, err
		}
		if isWicket {
			count++
			continue
		}
		break
	}
	return count, nil
}
