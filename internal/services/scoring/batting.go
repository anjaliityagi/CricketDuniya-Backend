package scoring

import "CricketDuniya-Backend/internal/dto"

func batting(req dto.BallRequest, ctx scoreContext) int {
	points := 0

	// runs
	switch req.RunsOffBat {
	case 1:
		points += 1
	case 2:
		points += 3
	case 4:
		points += 5
	case 6:
		points += 7
	}

	if req.IsBoundaryFour {
		n := ctx.prevConsecutiveBatterFours + 1
		if n >= 2 {
			points += 5 + (n * n)
		}
	}
	if req.IsBoundarySix {
		n := ctx.prevConsecutiveBatterSixes + 1
		if n >= 2 {
			points += 7 + (n * n)
		}
	}

	if req.IsDotBall {
		points -= 1
	}

	if req.IsWicket {
		points -= 3

		afterRuns := ctx.batterRunsBefore + req.RunsOffBat
		afterBalls := ctx.batterBallsBefore
		if req.BallType != "wide" && req.BallType != "no_ball" {
			afterBalls++
		}
		if afterRuns == 0 {
			points -= 5
			if afterBalls == 1 {
				points -= 7
			}
		}
	}

	afterRuns := ctx.batterRunsBefore + req.RunsOffBat
	if ctx.batterRunsBefore < 25 && afterRuns >= 25 {
		points += 15
	}
	if ctx.batterRunsBefore < 50 && afterRuns >= 50 {
		points += 30
	}
	if ctx.batterRunsBefore < 100 && afterRuns >= 100 {
		points += 50
	}

	afterBalls := ctx.batterBallsBefore
	if req.BallType != "wide" && req.BallType != "no_ball" {
		afterBalls++
	}
	beforeSR := ctx.batterStrikeRateBonusBefore
	afterSR := strikeRateBonus(afterRuns, afterBalls)
	points += (afterSR - beforeSR)

	return points
}

func strikeRateBonus(runs int, balls int) int {
	if balls < 5 || balls == 0 {
		return 0
	}
	sr := float64(runs) * 100.0 / float64(balls)
	switch {
	case sr < 20:
		return -5
	case sr < 40:
		return -3
	case sr < 60:
		return -2
	case sr < 80:
		return -1
	case sr < 100:
		return 0
	case sr < 140:
		return 1
	case sr < 160:
		return 2
	case sr < 180:
		return 3
	case sr < 200:
		return 5
	case sr < 220:
		return 8
	case sr < 240:
		return 13
	case sr < 260:
		return 21
	case sr < 280:
		return 34
	default:
		return 34
	}
}
