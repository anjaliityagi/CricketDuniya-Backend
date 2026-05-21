package scoring

import "CricketDuniya-Backend/internal/dto"

func bowling(req dto.BallRequest, ctx scoreContext) int {
	points := 0

	if req.IsDotBall {
		points += 2
	}

	if req.IsWicket {
		points += 5
		if req.DismissalType == "bowled" {
			points += 3
		}
		n := ctx.prevConsecutiveBowlerWicketsInOver + 1
		if n >= 2 {
			points += 5 + (n * n)
		}
	}

	if req.IsBoundaryFour {
		points -= 2
		n := ctx.prevConsecutiveBowlerFours + 1
		if n >= 2 {
			points -= (2 + (n * n))
		}
	}

	if req.IsBoundarySix {
		points -= 3
		n := ctx.prevConsecutiveBowlerSixes + 1
		if n >= 2 {
			points -= (3 + (n * n))
		}
	}

	if req.Wides > 0 {
		points -= 1
		points -= (ctx.prevConsecutiveBowlerExtras + 1)
	}
	if req.NoBalls > 0 {
		points -= 1
		points -= (ctx.prevConsecutiveBowlerExtras + 1)
	}

	legalBall := req.BallType != "wide" && req.BallType != "no_ball"
	overLegalAfter := ctx.overLegalBallsBefore
	if legalBall {
		overLegalAfter++
	}
	overRunsAfter := ctx.overRunsConcededBefore + req.TotalRuns - req.Byes - req.LegByes
	overBoundariesAfter := ctx.overBoundariesBefore
	if req.IsBoundaryFour || req.IsBoundarySix {
		overBoundariesAfter++
	}
	if overLegalAfter == 6 {
		if overRunsAfter == 0 {
			points += 20
		}
		if overBoundariesAfter == 0 {
			points += 5
		}

		ballsAfter := ctx.bowlerLegalBallsBefore
		if legalBall {
			ballsAfter++
		}
		runsAfter := ctx.bowlerRunsConcededBefore + req.TotalRuns - req.Byes - req.LegByes
		beforeEco := ctx.bowlerEconomyBonusBefore
		afterEco := economyBonus(runsAfter, ballsAfter)
		points += (afterEco - beforeEco)
	}

	return points
}

func economyBonus(runsConceded int, legalBalls int) int {
	if legalBalls < 6 {
		return 0
	}
	econ := float64(runsConceded) * 6.0 / float64(legalBalls)
	switch {
	case econ < 4:
		return 15
	case econ < 6:
		return 10
	case econ < 8:
		return 5
	case econ < 10:
		return 0
	case econ < 12:
		return -5
	default:
		return -10
	}
}
