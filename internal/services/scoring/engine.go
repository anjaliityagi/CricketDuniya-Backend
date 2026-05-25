package scoring

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
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

		t, err := repositories.GetBatterTotals(req.MatchID.String(), req.InningsID.String(), req.StrikerID.String())
		if err != nil {
			return ctx, err
		}
		ctx.batterRunsBefore, ctx.batterBallsBefore = t.RunsBefore, t.BallsBefore
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

		t, err := repositories.GetBowlerTotals(req.MatchID.String(), req.InningsID.String(), req.BowlerID.String())
		if err != nil {
			return ctx, err
		}
		ctx.bowlerLegalBallsBefore, ctx.bowlerRunsConcededBefore = t.LegalBallsBefore, t.RunsConcededBefore
		ctx.bowlerEconomyBonusBefore = economyBonus(ctx.bowlerRunsConcededBefore, ctx.bowlerLegalBallsBefore)

		t2, err := repositories.GetBowlerOverTotals(req.MatchID.String(), req.InningsID.String(), req.BowlerID.String(), req.BallNo)
		if err != nil {
			return ctx, err
		}
		ctx.overLegalBallsBefore, ctx.overRunsConcededBefore, ctx.overBoundariesBefore = t2.OverLegalBefore, t2.OverRunsBefore, t2.OverBoundariesBefore
	}

	return ctx, nil
}

func countConsecutiveBatterBoundary(req dto.BallRequest, boundaryType string) (int, error) {
	items, err := repositories.ListRecentRunsOffBatForStriker(req.MatchID.String(), req.InningsID.String(), req.StrikerID.String(), 24)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, runsOffBat := range items {
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
	items, err := repositories.ListRecentRunsOffBatForBowler(req.MatchID.String(), req.InningsID.String(), req.BowlerID.String(), 24)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, runsOffBat := range items {
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
	items, err := repositories.ListRecentBallTypesForBowler(req.MatchID.String(), req.InningsID.String(), req.BowlerID.String(), 24)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, t := range items {
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
	items, err := repositories.ListRecentWicketsForBowlerInOver(req.MatchID.String(), req.InningsID.String(), req.BowlerID.String(), req.BallNo, 24)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, isWicket := range items {
		if isWicket {
			count++
			continue
		}
		break
	}
	return count, nil
}
