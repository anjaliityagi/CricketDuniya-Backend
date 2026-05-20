package scoring

import "CricketDuniya-Backend/internal/dto"

func bowling(req dto.BallRequest) int {
	points := 0

	if req.IsDotBall {
		points += 2
	}

	if req.IsWicket {
		points += 5
	}

	if req.IsBoundaryFour {
		points -= 2
	}

	if req.IsBoundarySix {
		points -= 3
	}

	if req.Wides > 0 || req.NoBalls > 0 {
		points -= 1
	}

	return points
}
